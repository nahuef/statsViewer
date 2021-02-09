const hide = "none";
const show = "block";
const showStr = "Show chart";
const hideStr = "Hide chart";
const noChartStr = "No chart, less than 2 datapoints or 3 challenges.";

function toggleChart(name, action) {
    const c = document.getElementById("chart"+name);
    const btn = document.getElementById("btn"+name);

    // If no chart
    if (!c.children.length) {
        const scenario = document.getElementById("name"+name);
        scenario.title = noChartStr;
        c.style.display = hide;
        btn.style.display = hide;
        return
    }

    if (action && action === hide || !action && c.style.display === show) {
        c.style.display = hide;
        btn.innerText = showStr;
    } else {
        c.style.display = show;
        btn.innerText = hideStr;
    }
}

function toggleAll(action) {
    Array.from(document.getElementsByClassName("charts")).forEach(c =>  toggleChart(c.id.replace('chart', ''), action));
}
