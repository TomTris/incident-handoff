import {login} from "./api.js";

const errorDiv = document.getElementById("error-message") as HTMLDivElement;
const form = document.getElementById("login-form") as HTMLFormElement;

form?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const username = (document.getElementById("username") as HTMLInputElement).value
    const password = (document.getElementById("password") as HTMLInputElement).value

    try {
        await login(username, password);
        window.location.href = "/"
    } catch (err) {
        errorDiv.textContent = (err as Error).message;
        errorDiv.style.display = "block";
    }
});