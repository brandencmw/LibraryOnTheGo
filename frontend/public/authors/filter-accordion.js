const accordion = document.getElementById("filterAccordionHeader")

accordion.addEventListener("click", event => {
    if (event.target === accordion) {
        const content = document.getElementById("filterAccordionContent")
        content.style.display = content.style.display === 'none' ? content.style.display = 'block' : 'none';
    }
})

const sortDropdown = document.getElementById("sortDropdown")
sortDropdown.addEventListener("change", event => {
    const currentURL = new URL(window.location.href)
    sort = event.target.value
    console.log(sort)
    currentURL.searchParams.set("order", sort)

    window.location.href = currentURL.href
})

const filterButtons = document.getElementsByClassName("alphabet-button")
for (let i = 0; i < filterButtons.length; i++) {
    filterButtons[i].addEventListener("click", event => {
        const letter = event.target.id.slice(-1)
        const currentURL = new URL(window.location.href)
        currentURL.searchParams.set("name", letter)

        window.location.href = currentURL.href
    })
}