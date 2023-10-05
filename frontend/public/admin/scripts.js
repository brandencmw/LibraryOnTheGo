function userAuthorized() {
    let url = "https://localhost:8080/auth"
    console.log(document.cookie)
    return fetch(url, {
        method: "POST",
        credentials: "include"
    }).then(response => {
        console.log(response)
        if (!response.ok) {
            throw new Error()
        }
        return response.json
    }).then(data => {
        console.log(data)
        return data
    }).catch(_ => {
        throw new Error("Unauthorized user")
    })
}

userAuthorized()
.then(response => {
    console.log(response);
}).catch(error => {
    console.error(error.message)
    // window.location.href = "/login"
})