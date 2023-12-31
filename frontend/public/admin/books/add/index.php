<?php
require '../../auth.php';
if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Add Book</title>
</head>
<body>
    <h2>Search for Books on Google Books</h2>
    <form id="googleBookSearchForm">
        <label for="google-book-title">Title</label>
        <input id="googleBookTitle" name="google-book-title" type="text">
        <label for="google-book-author">Author</label>
        <input id="googleBookAuthor" name="google-book-author" type="text">
        <button type="submit">Search for a Book</button>
    </form>

    <form id="addBookForm">
        <label for="add-book-title">Title</label>
        <input id="addBookTitle" name="add-book-title" type="text">
        <label for="add-book-authors">Authors</label>
        <input id="addBookAuthors" name="add-book-authors" type="text">
        <label for="add-book-publish-date">Publish Date</label>
        <input id="addBookPublishDate" name="add-book-publish-date" type="text">
        <label for="add-book-page-count">Page Count</label>
        <input id="addBookPageCount" name="add-book-page-count" type="text">
        <label for="add-book-categories">Categories</label>
        <input id="addBookCategories" name="add-book-categories" type="text">
        <label for="add-book-description">Description</label>
        <textarea id="addBookDescription" name="add-book-description"cols="30" rows="10"></textarea>
        <label for="add-book-cover-photo">Cover Photo</label>
        <input id="addBookCover" type="file">
        <button type="submit">Add Book</button>
    </form>

    <h3>Google Results</h3>
    <table id="googleSearchResults">
        <thead>
            <th>Title</th>
            <th>Authors</th>
            <th>Publish Year</th>
            <th>Description</th>
        </thead>
        <tbody></tbody>
    </table>

</body>
<script src="./search-books.js"></script>
<script src="./add-book.js"></script>
<script src="../../send-data.js"></script>
</html>

<?php
} else {
    header('../../../login/');
}
?>