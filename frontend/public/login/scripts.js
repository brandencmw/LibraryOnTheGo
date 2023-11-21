async function requestLogin(requestData) {
    const url = "https://localhost:8080/auth/login";
    const options = {
        method: "POST", 
        credentials: "include", 
        body: JSON.stringify(requestData)
    }
    const response = await fetch(url, options)
    if (!response.ok) {
        throw new Error("Response not ok")
    }
    const responseJSON = await response.json()
    return responseJSON
}


function onSubmit(event) {
    event.preventDefault();
    const request = {
        username: form.elements["username"].value,
        password: form.elements["password"].value,
    };
    requestLogin(request)
        .then(window.location.href = "/admin")
        .catch(err => console.error(err))
}

const form = document.getElementById("loginForm");
form.addEventListener("submit", (event) => onSubmit(event));