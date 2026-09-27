(function () {
    var input = document.getElementById("searchQuery_small");
    if (!input) return;
    var form = input.closest ? input.closest("form") : null;
    var resultsList = document.getElementById("results_small");
    var suggestions = [];
    var activeIndex = -1;

    function escapeRegExp(s) {
        return String(s).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    }

    // Enter always lands on the paginated results page (/search?q=...).
    // Clicking (or arrow-picking + Enter on) a suggestion opens that product.
    function goToSearchPage(q) {
        q = (q || "").trim();
        if (!q) return;
        window.location.href = "/search?q=" + encodeURIComponent(q);
    }

    function setActive(idx) {
        activeIndex = idx;
        if (!resultsList) return;
        var items = resultsList.querySelectorAll("li");
        items.forEach(function (li, i) {
            if (i === idx) {
                li.classList.add("active");
                li.style.backgroundColor = "#f0f0f0";
            } else {
                li.classList.remove("active");
                li.style.backgroundColor = "";
            }
        });
    }

    function render(data, query) {
        if (!resultsList) return;
        resultsList.innerHTML = "";
        suggestions = [];
        activeIndex = -1;

        if (data.results && Array.isArray(data.results) && data.results.length > 0) {
            suggestions = data.results;
            var rx = new RegExp(escapeRegExp(query), "gi");
            data.results.forEach(function (result, idx) {
                var highlighted = String(result.title).replace(rx, function (m) {
                    return '<span class="highlight_search">' + m + "</span>";
                });
                var li = document.createElement("li");
                var a = document.createElement("a");
                a.href = result.link;
                a.innerHTML = highlighted;
                a.addEventListener("mouseenter", function () { setActive(idx); });
                li.appendChild(a);
                resultsList.appendChild(li);
            });
        } else if ((query || "").trim().length > 0) {
            resultsList.innerHTML = "<li>نتیجه‌ای یافت نشد. برای دیدن همه نتایج اینتر را بزنید.</li>";
        }
    }

    function resolveEnter() {
        if (activeIndex >= 0 && suggestions[activeIndex]) {
            window.location.href = suggestions[activeIndex].link;
        } else {
            goToSearchPage(input.value);
        }
    }

    var debounce;
    input.addEventListener("input", function () {
        var query = this.value;
        clearTimeout(debounce);
        if (!query || query.trim().length === 0) {
            if (resultsList) resultsList.innerHTML = "";
            suggestions = [];
            activeIndex = -1;
            return;
        }
        debounce = setTimeout(function () {
            fetch("/tsearch?q=" + encodeURIComponent(query))
                .then(function (response) { return response.json(); })
                .then(function (data) {
                    if (input.value !== query) return; // stale response
                    render(data, query);
                })
                .catch(function (error) {
                    console.error("Error:", error);
                });
        }, 200);
    });

    input.addEventListener("keydown", function (e) {
        if (e.key === "ArrowDown") {
            if (!suggestions.length) return;
            e.preventDefault();
            var next = (activeIndex + 1) % suggestions.length;
            setActive(next);
            var itemsDown = resultsList.querySelectorAll("li");
            if (itemsDown[next]) itemsDown[next].scrollIntoView({ block: "nearest" });
        } else if (e.key === "ArrowUp") {
            if (!suggestions.length) return;
            e.preventDefault();
            var prev = (activeIndex - 1 + suggestions.length) % suggestions.length;
            setActive(prev);
            var itemsUp = resultsList.querySelectorAll("li");
            if (itemsUp[prev]) itemsUp[prev].scrollIntoView({ block: "nearest" });
        } else if (e.key === "Enter") {
            e.preventDefault();
            resolveEnter();
        } else if (e.key === "Escape") {
            if (resultsList) resultsList.innerHTML = "";
            activeIndex = -1;
        }
    });

    if (form) {
        form.addEventListener("submit", function (e) {
            e.preventDefault();
            resolveEnter();
        });
    }
})();
