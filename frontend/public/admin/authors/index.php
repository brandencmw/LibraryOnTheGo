<?php

    require './get_authors.php';
    require '../auth.php';

    if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
        // Get result into array
        $authors = getAuthors();
?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Author Operations</title>
</head>
<body>
    <h2>Authors in the System</h2>
    <table id="authorTable">
        <thead>
            <th>ID</th>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Actions</th>
        </thead>
        <tbody>
        <?php
            // Cycle through array and add rows to table
            foreach ($authors as $author) {
        ?>
            <tr id="author-<?php echo $author['id'] ?>">
                <td><?php echo $author['id'] ?></td>
                <td><?php echo $author['firstName'] ?></td>
                <td><?php echo $author['lastName'] ?></td>
                <td>
                    <button id="update-<?php echo $author['id'] ?>">Update</button>
                    <button id="delete-<?php echo $author['id'] ?>">Delete</button>
                </td>
            </tr>
        <?php
            }
        ?>
        </tbody>
    </table>
    <h2>Create New Author</h2>
    <form method="post" id="addAuthorForm" enctype="multipart/form-data">
        <label for="authorFirstName">Author's First Name</label>
        <input type="text" name="authorFirstName" id="authorFirstName">

        <label for="authorLastName">Author's Last Name</label>
        <input type="text" name="authorLastName" id="authorLastName">

        <label for="authorBio">Author's Bio</label>
        <textarea name="authorBio" id="authorBio" cols="30" rows="10"></textarea>

        <label for="authorHeadshot">Upload Headshot</label>
        <input type="file" name="authorHeadshot" id="authorHeadshot">
        <button type="submit">Submit</button>
    </form>
</body>
<script src="add-author.js"></script>
<script src="button-listeners.js"></script>
<script src="../send-data.js"></script>
</html>

<?php
    } else {
        header('../../../login/');
    }
?>