function onSubmit(event) {
    event.preventDefault();
    const form = document.getElementById("loginForm");
    const url = "http://localhost:8080/login";
    const username = form.elements["username"].value;
    const password = form.elements["password"].value;

    const request = {
        username,
        password,
    };

    fetch(url, {
        method: "POST",
        headers: {
            'Content-Type': "application/json"
        },
        body: JSON.stringify(request)
    }).then(response => {
        console.log(response);
        console.log(response.headers)
        if (!response.ok) {
            throw new Error("Response not ok")
        }
        return response.json();
    }).then(data => {
        console.log(data);
        //window.location.href = "/admin"
    }).catch(error => {
        console.log(error);
    })

}

document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("loginForm");
    form.addEventListener("submit", (event) => onSubmit(event));
});