document.getElementById("searchQuery").addEventListener("input", function () {
    const query = this.value;

    if (query.length > 0) {
        fetch(`/tsearch?q=${query}`)
            .then(response => response.json())
            .then(data => {
                const resultsList = document.getElementById("results");
                resultsList.innerHTML = '';

                if (data.results && Array.isArray(data.results)) {
                    data.results.forEach(result => {
                        const highlightedTitle = result.title.replace(new RegExp(query, 'gi'), (match) => {
                            return `<span class="highlight_search">${match}</span>`;
                        });

                        const li = document.createElement("li");
                        li.innerHTML = `<a  href="${result.link}">${highlightedTitle}</a>`;
                        resultsList.appendChild(li);
                    });
                } else {
                    resultsList.innerHTML = `<li>نتیجه‌ای یافت نشد.</li>`;
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }
});