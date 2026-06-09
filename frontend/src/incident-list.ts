import { createIncident, loadIncidents, logout } from "./api.js"
import { domDiv, domForm, domUl, inputValue } from "./dom.js"
import type { Incident, CreateIncidentRequest, Severity } from "./types.js"

var incidents : Incident[] = []
const ul = domUl("incident-list")

function render() {
    console.log(123)
    ul.innerHTML = ""
    console.log(2)
    incidents.forEach(inc => {
        console.log(3)
        const date = new Date(inc.created_at);
        const formatted = date.toISOString().slice(5, 16).replace('T', ' ');
        const li = document.createElement("li")
        li.innerHTML = `<span class=${inc.severity}>${inc.severity}</span>
        <h3>${inc.title}</h3>
        <p>${formatted}</p>
        <p>Status: ${inc.status}</p>`
        ul.appendChild(li)
    })
    console.log(4)
}

async function refreshIncidents() {
    console.log("12344422222")
    incidents = await loadIncidents()
    await render()
}
    
const form = domForm("add-incident-form")
form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const errorDiv = domDiv("error-message")

    function isSeverity(v: string) : v is Severity {
        return v === "SEV1" || v === "SEV2" || v === "SEV3"
    }
    
    const severityRaw: string = inputValue("incident-severity")
    if (!isSeverity(severityRaw)) {
        errorDiv.textContent = "Severity must be SEV1, SEV2 or SEV3"
        errorDiv.style = "block"
        return
    }

    const input : CreateIncidentRequest = {
        title: inputValue("incident-title"),
        service: inputValue("incident-service"),
        severity: severityRaw,
    }
    try {
        await createIncident(input)
        await refreshIncidents()
    } catch (err) {
        errorDiv.textContent = (err as Error).message;
        errorDiv.style.display = "block";
    }
})

const logoutBtn = document.getElementById("logout-btn") as HTMLButtonElement;
logoutBtn.addEventListener("click", async (event) => {
    event.preventDefault();
    try {
        await logout()
    }
    catch (err) {
        const name: string = (err as Error).name
        const message: string = (err as Error).message
        alert(`${name}: ${message}`)
    } finally {
        window.location.href = "/"
    }
})

console.log("123")
refreshIncidents()
console.log("123444")