const Controller = {
  results: [],
  hasMore: false,

  search: (ev) => {
    ev.preventDefault();
    const data = getFormData();
    fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.results = results;
        Controller.hasMore = response.headers.get("X-Has-More") === "true";
        Controller.updateTable();
      });
    });
  },

  loadMore: () => {
    if (!Controller.hasMore) return;
    
    const data = getFormData();
    const offset = Controller.results.length;
    fetch(`/search?q=${data.query}&offset=${offset}`).then((response) => {
      response.json().then((results) => {
        Controller.results.push(...results);
        Controller.hasMore = response.headers.get("X-Has-More") === "true";
        Controller.updateTable();
      });
    });
  },

  updateTable: () => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of Controller.results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows;
  },
};

function getFormData() {
  const form = document.getElementById("form");
  return Object.fromEntries(new FormData(form));
}

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", Controller.loadMore);
