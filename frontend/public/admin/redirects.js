
authorButton = document.getElementById("authorPageButton");

authorButton.addEventListener("click", _ => {
    window.location.href = "/admin/authors"
})

bookButton = document.getElementById("bookPageButton");

bookButton.addEventListener("click", _ => {
    window.location.href = "/admin/books"
})