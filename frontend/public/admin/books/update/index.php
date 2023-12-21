<?php declare(strict_types=1);
require '../../auth.php';

function getBookToUpdate(string $id) {
    $get_book_endpoint = 'https://server/books?id=' . $id;

    $ch = curl_init($get_book_endpoint);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'GET');
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');

    $json_response = curl_exec($ch);
    if (curl_errno($ch)) {
        echo ''. curl_error($ch);
        return null;
    }
    curl_close($ch);

    return json_decode($json_response, true)['book'];
}

function formatAuthorList(array $authors) {
    $names = [];
    foreach ($authors as $author) {
        $names[] = $author['firstName'] . ' ' . $author['lastName'];
    }

    return implode(',', $names);
}

if (isset($_COOKIE['Authorization']) && authenticated($_COOKIE['Authorization'])) {
    if (isset($_GET['id'])) {
        $book = getBookToUpdate($_GET['id']);
        if ($book == null) {
            header('Location: ../');
        }
?>

        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Update Book</title>
        </head>
        <body>
            <form id="updateBookForm">
                <label for="bookTitle">Title</label>
                <input name="bookTitle" id="bookTitle" type="text" value="<?php echo isset($book['title']) ? $book['title'] : ''; ?>" disabled>

                <label for="bookAuthors">Authors</label>
                <input name="bookAuthors" id="bookAuthors" type="text" value="<?php echo isset($book['authors']) ? formatAuthorList($book['authors']) : ''; ?>" disabled>

                <label for="bookPublishDate">Publish Date</label>
                <input name="bookPublishDate" id="bookPublishDate" type="text" value="<?php echo isset($book['publishDate']) ? $book['publishDate'] : ''; ?>" disabled>

                <label for="bookPageCount">Page Count</label>
                <input name="bookPageCount" id="bookPageCount" type="text" value="<?php echo isset($book['pageCount']) ? $book['pageCount'] : ''; ?>" disabled>

                <label for="bookCategories">Categories</label>
                <input name="bookCategories" id="bookCategories" type="text" value="<?php echo isset($book['categories']) ? implode(',', $book['categories']) : ''; ?>" disabled>

                <label for="bookSynopsis">Synopsis</label>
                <textarea name="bookSynopsis" id="bookSynopsis" cols="30" rows="10"><? echo isset($book['synopsis']) ? $book['synopsis'] : ''; ?></textarea>

                <label for="bookCover">Cover Image</label>
                <input name="bookCover" id="bookCover" type="file">
                <img id="coverPreview" src="<?php echo isset($book['coverKey']) ? 'https://library-pictures.s3.amazonaws.com/' . $book['coverKey'] : '#' ?>">
                <button type="submit">Submit</button>
            </form>
            <button id="backButton">Back</button>
        </body>
        <script src="../../send-data.js"></script>
        <script src="./update-book.js"></script>
        </html>

<?php
    } else {
        header('Location: ../');
    }
} else {
    header('Location: ../../../login/');
}
?>