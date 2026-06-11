import { isAuthenticated } from "./api.js";
async function checkAuth() {
    try {
        const authed = await isAuthenticated();
        window.location.href = authed ? "/incident-list.html" : "/login.html"
    } catch {
        alert("Unexpected Error");
    }
}
checkAuth()