package analyzer

import (
	"encoding/json"
	"os"
)

func GenerateHTMLReport(output string, reports []ModuleReport) error {

	data, _ := json.MarshalIndent(reports, "", "  ")

	html := `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>CI3 HMVC Analyzer</title>

<style>
body {
	margin: 0;
	font-family: Arial, sans-serif;
	display: flex;
	height: 100vh;
}

#sidebar {
	width: 300px;
	background: #1e1e2f;
	color: #fff;
	overflow-y: auto;
	padding: 10px;
}

#search {
	width: 100%;
	padding: 8px;
	border-radius: 4px;
	border: none;
	margin-bottom: 10px;
	font-size: 14px;
}

#allIssuesBtn {
	position: fixed;
	bottom: 20px;
	right: 20px;
	background: #dc3545;
	color: #fff;
	border: none;
	padding: 12px 18px;
	font-size: 14px;
	border-radius: 30px;
	cursor: pointer;
	box-shadow: 0 4px 10px rgba(0,0,0,0.3);
	z-index: 9999;
}

#allIssuesBtn:hover {
	background: #bd2130;
}

.module-title {
	cursor: pointer;
	font-size: 15px;
	margin: 10px 0;
	padding: 6px;
	background: #29293d;
	border-radius: 4px;
	display: flex;
	justify-content: space-between;
}

.class-list {
	margin-left: 10px;
	display: none;
}

.class-item {
	cursor: pointer;
	padding: 6px 10px;
	margin: 4px 0;
	background: #2b2b3d;
	border-radius: 4px;
	font-size: 14px;
}

#content {
	flex: 1;
	padding: 20px;
	overflow-y: auto;
	background: #f7f7f7;
}

.card {
	background: #fff;
	padding: 20px;
	border-radius: 6px;
	box-shadow: 0 2px 6px rgba(0,0,0,0.1);
}

.method {
	padding: 6px 10px;
	background: #efefef;
	margin: 5px 0;
	border-radius: 4px;
	font-family: monospace;
}

.badge {
	padding: 2px 6px;
	font-size: 12px;
	background: #007bff;
	color: #fff;
	border-radius: 4px;
	margin-left: 6px;
}

.warning {
	padding: 10px;
	margin: 8px 0;
	border-radius: 6px;
	color: #fff;
	font-size: 14px;
}

.warning pre {
	background: rgba(0,0,0,0.2);
	padding: 6px;
	border-radius: 4px;
	margin-top: 6px;
	white-space: pre-wrap;
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fill, 220px);
  gap: 16px;
}

.card {
  position: relative;
  padding: 16px 18px;
  border-radius: 12px;
  min-height: 110px;
  background: linear-gradient(135deg, #f9fafb, #ffffff);
  box-shadow: 0 6px 16px rgba(0,0,0,0.08);
}

.card h3 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #111827;
}

.card p {
  margin-top: 4px;
  font-size: 12px;
  color: #4b5563;
  text-transform: uppercase;
}

.card::after {
  content: attr(data-icon);
  position: absolute;
  right: 14px;
  bottom: 10px;
  font-size: 32px;
  opacity: 0.15;
}

.card:hover {
  transform: translateY(-4px);
  transition: 0.2s ease;
}

#footer {
	position: fixed;
	bottom: 0;
	left: 300px;
	right: 0;
	height: 20px;
	// background: #c2c2c5ff;
	color: black;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 13px;
	z-index: 999;
}


</style>
</head>

<body>

<div id="sidebar">
	<input id="search" placeholder="Search module / class / method" />
	<div id="moduleList"></div>
</div>

<div id="content">
	<div class="card">
		<h2>CI3 HMVC Analyzer</h2>
		<p>Select a class from the sidebar.</p>
	</div>
</div>

<button id="allIssuesBtn">🚨 View All Security Issues</button>

<script>
const reports = ` + string(data) + `;

const moduleList = document.getElementById("moduleList");
const content = document.getElementById("content");
const search = document.getElementById("search");
const allBtn = document.getElementById("allIssuesBtn");

 function renderSidebar(filter) {
	moduleList.innerHTML = "";
	filter = (filter || "").toLowerCase();

	reports.forEach(function(module) {

		const moduleName = module.Module.toLowerCase();
		let moduleMatched = moduleName.includes(filter);

		const moduleDiv = document.createElement("div");

		const title = document.createElement("div");
		title.className = "module-title";
		title.innerHTML = "<span>" + module.Module + "</span><span>▸</span>";

		const classList = document.createElement("div");
		classList.className = "class-list";

		let classMatched = false;

		if (module.Files && Array.isArray(module.Files)) {
			module.Files.forEach(function(file) {

				const className = file.ClassName.toLowerCase();
				let methodMatched = false;

				if (file.Methods) {
					methodMatched = file.Methods.some(m =>
						m.toLowerCase().includes(filter)
					);
				}

				if (
					className.includes(filter) ||
					methodMatched ||
					moduleMatched
				) {
					classMatched = true;

					const classDiv = document.createElement("div");
					classDiv.className = "class-item";
					classDiv.textContent = file.ClassName;
					classDiv.onclick = function () {
						showClass(module, file);
					};
					classList.appendChild(classDiv);
				}
			});
		}

		if (!moduleMatched && !classMatched && filter !== "") return;

		title.onclick = function () {
			const open = classList.style.display === "block";
			classList.style.display = open ? "none" : "block";
			title.lastChild.textContent = open ? "▸" : "▾";
		};

		moduleDiv.appendChild(title);
		moduleDiv.appendChild(classList);
		moduleList.appendChild(moduleDiv);
	});
}


function capitalizeFirst(str) {
  if (!str) return str;
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function showClass(module, file) {
	let html = "<div class='card'>";
	html += "<h2>" + capitalizeFirst(file.Folder) + "</h2>";
	html += "<h2>" + file.ClassName + "</h2>";
	html += "<p><strong>Module:</strong> " + module.Module + "</p>";
	html += "<p><strong>File:</strong> <code>" + file.File + "</code></p>";

	html += "<h3>Methods <span class='badge'>" + file.Methods.length + "</span></h3>";
	file.Methods.forEach(function(m) {
		html += "<div class='method'>" + m + "()</div>";
	});

	html += "<h3>Security Warnings</h3>";

	if (file.Warnings && file.Warnings.length > 0) {
		file.Warnings.forEach(function (w) {
		html += renderWarning(w);
	});
	} else {
		html += "<p style='color:green'>✅ No security issues detected.</p>";
	}

	html += "</div>";
	content.innerHTML = html;
}

function renderWarning(w) {
	let bg = "#6c757d";
	if (w.Level === "HIGH") bg = "#dc3545";
	else if (w.Level === "MEDIUM") bg = "#ffc107";
	else if (w.Level === "LOW") bg = "#17a2b8";

	return (
		"<div class='warning' style='background:" + bg + "'>" +
		"<strong>" + w.Level + "</strong> - " + w.Message + "<br>" +
		"<small>📄 " + w.File + " : Line " + w.Line + "</small>" +
		"<pre><code>" + (w.Snippet || "") + "</code></pre>" +
		"</div>"
	);
}

function showAllIssues() {
	let html = "<div class='card'><h1>🚨 Project Security Issues</h1>";
	let found = false;

	reports.forEach(function(module) {
	
		if(module.Files && Array.isArray(module.Files)){

		module.Files.forEach(function(file) {
			if (!file.Warnings || file.Warnings.length === 0) return;

				found = true;
				html += "<h3>" + module.Module + " / " + file.ClassName + "</h3>";
				file.Warnings.forEach(function(w) {
					html += renderWarning(w);
				});
			});
		}
	});

	if (!found) {
		html += "<p style='color:green'>✅ No security issues found.</p>";
	}

	html += "</div>";
	content.innerHTML = html;
}

allBtn.onclick = showAllIssues;
renderSidebar("");

function createCard(value, label, type, icon) {
  var card = document.createElement("div");
  card.className = "card " + type;
  card.setAttribute("data-icon", icon);

  var h3 = document.createElement("h3");
  h3.innerText = value;

  var p = document.createElement("p");
  p.innerText = label;

  card.appendChild(h3);
  card.appendChild(p);

  return card;
}

function renderDashboard(reports) {
  var moduleCount = reports.length;
  var fileCount = 0;
  var classCount = 0;
  var methodCount = 0;
  var warningCount = 0;

  reports.forEach(function (module) {
    if (module.Files) {
      fileCount += module.Files.length;

      module.Files.forEach(function (file) {
        if (file.ClassName) classCount++;
        if (file.Methods) methodCount += file.Methods.length;
        if (file.Warnings) warningCount += file.Warnings.length;
      });
    }
  });

  var content = document.getElementById("content");
  content.innerHTML = "";

  var dashboard = document.createElement("div");
  dashboard.className = "dashboard";

  var title = document.createElement("h2");
  title.innerText = "CI3 HMVC Analyzer – Overview";
  dashboard.appendChild(title);

  var stats = document.createElement("div");
  stats.className = "stats";

  stats.appendChild(createCard(moduleCount, "Modules", "modules", "📦"));
  stats.appendChild(createCard(fileCount, "Files", "files", "📄"));
  stats.appendChild(createCard(classCount, "Classes", "classes", "🧠"));
  stats.appendChild(createCard(methodCount, "Methods", "methods", "🔧"));
  stats.appendChild(createCard(warningCount, "Warnings", "warning", "⚠️"));

  dashboard.appendChild(stats);
  content.appendChild(dashboard);
}

renderDashboard(reports);

search.addEventListener("input", function () {
	renderSidebar(this.value);
});

document.getElementById("year").innerText = new Date().getFullYear();

</script>
<footer id="footer">
	© <span id="year"></span> Vicky Chhetri
</footer>

</body>
</html>`

	return os.WriteFile(output, []byte(html), 0644)
}
