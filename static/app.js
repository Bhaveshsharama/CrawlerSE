let limit = 10;
let offset = 0;
let currentQuery = "";

const form = document.getElementById("searchForm");
const box = document.getElementById("searchBox");
const resultsDiv = document.getElementById("results");
const infoDiv = document.getElementById("info");
const pagination = document.getElementById("pagination");
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");

form.addEventListener("submit", function (e) {
    e.preventDefault();
    offset = 0;
    currentQuery = box.value.trim();
    if (currentQuery) runSearch();
});

prevBtn.addEventListener("click", function () {
    if (offset >= limit) {
        offset -= limit;
        runSearch();
    }
});

nextBtn.addEventListener("click", function () {
    offset += limit;
    runSearch();
});

async function runSearch() {
    infoDiv.textContent = "Searching...";
    resultsDiv.innerHTML = "";
    pagination.style.display = "none";

    try {
        const url = `/search?q=${encodeURIComponent(currentQuery)}&limit=${limit}&offset=${offset}`;
        const res = await fetch(url);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        renderResults(data);
    } catch (err) {
        infoDiv.textContent = "";
        resultsDiv.innerHTML = `<div class="empty-state">Search failed — is the server running?</div>`;
    }
}

function escapeHTML(str) {
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}

function renderResults(data) {
    resultsDiv.innerHTML = "";

    const page = Math.floor(offset / limit) + 1;
    const totalPages = Math.ceil(data.total / limit);

    if (data.total === 0) {
        infoDiv.textContent = `No results · ${data.time_ms} ms`;
        pagination.style.display = "none";
        resultsDiv.innerHTML = `<div class="empty-state">No results found for "${escapeHTML(currentQuery)}"</div>`;
        return;
    }

    infoDiv.textContent = `${data.total} results · page ${page}/${totalPages} · ${data.time_ms} ms`;

    data.results.forEach(r => {
        const div = document.createElement("div");
        div.className = "result";

        const link = document.createElement("a");
        link.href = r.url;
        link.target = "_blank";
        link.rel = "noopener";
        link.textContent = r.title || r.url;

        const urlLine = document.createElement("div");
        urlLine.className = "url-display";
        urlLine.textContent = r.url;

        const snippet = document.createElement("div");
        snippet.className = "snippet";
        snippet.textContent = r.snippets || "";

        const score = document.createElement("div");
        score.className = "score";
        score.textContent = `score: ${r.score.toFixed(2)}`;

        div.appendChild(link);
        div.appendChild(urlLine);
        div.appendChild(snippet);
        div.appendChild(score);
        resultsDiv.appendChild(div);
    });

    prevBtn.disabled = offset === 0;
    nextBtn.disabled = offset + limit >= data.total;
    pagination.style.display = "flex";
}

window.onload = function () {
    const params = new URLSearchParams(window.location.search);
    const q = params.get("q");
    const off = params.get("offset");

    if (q) {
        currentQuery = q;
        box.value = q;
        offset = off ? parseInt(off) : 0;
        runSearch();
    }
};