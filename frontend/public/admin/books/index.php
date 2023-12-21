<?
require './get_books.php';
require '../auth.php';

if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
    $books = getAllBooks()
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Book Operations</title>
</head>
<body>
    <button id="addBookButton">Add a new book</button>
    <table id="bookTable">
        <thead>
            <th>ID</th>
            <th>Title</th>
            <th>Authors</th>
            <th>Actions</th>
        </thead>
        <tbody>
            <?php foreach ($books as $book) { ?>
            <tr id="book-<?php echo $book['id']; ?>">
                <td><?php echo $book['id'] ?></td>
                <td><?php echo $book['title']; ?></td>
                <?php 
                $authorList = array();
                foreach ($book['authors'] as $author) { 
                    $authorList[] = $author['firstName'] . ' ' . $author['lastName'];
                }
                ?>
                <td><?php echo join($authorList); ?></td>
                <td>
                    <button id="update-<?php echo $book['id'] ?>">Update</button>
                    <button id="delete-<?php echo $book['id'] ?>">Delete</button>
                </td>
            </tr>
            <?php } ?>
        </tbody>
    </table>
</body>
<script src="./redirect.js"></script>
<script src="./button-listeners.js"></script>
</html>

<?php
} else {
    header('Location: ../../login/');
}
?>