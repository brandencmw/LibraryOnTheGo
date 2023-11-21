async function requestAuth() {
    const url = "https://localhost:8080/auth/"

    const response = await fetch(url, {method: "POST", credentials: "include"})
    if (!response.ok) throw new Error("Unauthorized User");
    const responseJSON = response.json()
    return responseJSON
}

requestAuth()
    .catch(err => {
        alert(err);
        window.location.href = "/login"
    })