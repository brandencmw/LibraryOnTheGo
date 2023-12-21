<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- Stylesheets -->
    <link rel="stylesheet" href="styles.css">

    <!-- Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Merriweather&family=Poppins&display=swap" rel="stylesheet">
    <title>Library</title>
</head>
<body>
    <div class="main-box">
        <header>
            <h1 class="main-heading">Branden's Library</h1>
            <div class="header-links">
                <a href="/authors">See all authors</a>
                <a href="/books">See all books</a>
            </div>
        </header>

        <div class="content">
            <form action="" method="GET" class="search-form">
                <label class="search-label" for="search-terms">Search</label>
                <div class="search-bar">
                    <input class="search-terms" type="text" name="search-terms" id="search-terms">
                    <button class="search-btn" type="submit">
                        <svg id="search-icon" class="search-icon" viewBox="0 0 24 24">
                            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                            <path d="M0 0h24v24H0z" fill="none"/>
                        </svg>
                    </button>
                </div>
            </form>
        </div>

        <footer>
            <svg fill="#FFFFFF" xmlns="http://www.w3.org/2000/svg"  viewBox="0 0 50 50" width="30px" height="30px">    
                <path d="M41,4H9C6.24,4,4,6.24,4,9v32c0,2.76,2.24,5,5,5h32c2.76,0,5-2.24,5-5V9C46,6.24,43.76,4,41,4z M17,20v19h-6V20H17z M11,14.47c0-1.4,1.2-2.47,3-2.47s2.93,1.07,3,2.47c0,1.4-1.12,2.53-3,2.53C12.2,17,11,15.87,11,14.47z M39,39h-6c0,0,0-9.26,0-10 c0-2-1-4-3.5-4.04h-0.08C27,24.96,26,27.02,26,29c0,0.91,0,10,0,10h-6V20h6v2.56c0,0,1.93-2.56,5.81-2.56 c3.97,0,7.19,2.73,7.19,8.26V39z"/>
            </svg>
            <svg width="37" height="30" viewBox="0 0 68 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M68 0H0V48H68V0ZM17.8667 18.3359L16.6436 25.0229L15.1567 18.3359H12.9741L11.5093 25.0229L10.2861 18.3359H7.72998L10.0005 29H12.6592L14.0581 22.9209L15.4863 29H18.1523L20.4155 18.3359H17.8667ZM34.6538 18.3359L33.4307 25.0229L31.9438 18.3359H29.7612L28.2964 25.0229L27.0732 18.3359H24.5171L26.7876 29H29.4463L30.8452 22.9209L32.2734 29H34.9395L37.2026 18.3359H34.6538ZM51.4409 18.3359L50.2178 25.0229L48.731 18.3359H46.5483L45.0835 25.0229L43.8604 18.3359H41.3042L43.5747 29H46.2334L47.6323 22.9209L49.0605 29H51.7266L53.9897 18.3359H51.4409ZM60.0762 26.9565C59.8174 26.7222 59.4927 26.605 59.1021 26.605C58.7065 26.605 58.3794 26.7222 58.1206 26.9565C57.8667 27.1909 57.7397 27.4863 57.7397 27.8428C57.7397 28.1992 57.8667 28.4946 58.1206 28.729C58.3794 28.9634 58.7065 29.0806 59.1021 29.0806C59.4927 29.0806 59.8174 28.9658 60.0762 28.7363C60.335 28.502 60.4644 28.2041 60.4644 27.8428C60.4644 27.4814 60.335 27.186 60.0762 26.9565Z" fill="white"/>
            </svg>
            <svg fill="#FFFFFF" xmlns="http://www.w3.org/2000/svg"  viewBox="0 0 32 32" width="35px" height="35px"><path fill-rule="evenodd" d="M 16 4 C 9.371094 4 4 9.371094 4 16 C 4 21.300781 7.4375 25.800781 12.207031 27.386719 C 12.808594 27.496094 13.027344 27.128906 13.027344 26.808594 C 13.027344 26.523438 13.015625 25.769531 13.011719 24.769531 C 9.671875 25.492188 8.96875 23.160156 8.96875 23.160156 C 8.421875 21.773438 7.636719 21.402344 7.636719 21.402344 C 6.546875 20.660156 7.71875 20.675781 7.71875 20.675781 C 8.921875 20.761719 9.554688 21.910156 9.554688 21.910156 C 10.625 23.746094 12.363281 23.214844 13.046875 22.910156 C 13.15625 22.132813 13.46875 21.605469 13.808594 21.304688 C 11.144531 21.003906 8.34375 19.972656 8.34375 15.375 C 8.34375 14.0625 8.8125 12.992188 9.578125 12.152344 C 9.457031 11.851563 9.042969 10.628906 9.695313 8.976563 C 9.695313 8.976563 10.703125 8.65625 12.996094 10.207031 C 13.953125 9.941406 14.980469 9.808594 16 9.804688 C 17.019531 9.808594 18.046875 9.941406 19.003906 10.207031 C 21.296875 8.65625 22.300781 8.976563 22.300781 8.976563 C 22.957031 10.628906 22.546875 11.851563 22.421875 12.152344 C 23.191406 12.992188 23.652344 14.0625 23.652344 15.375 C 23.652344 19.984375 20.847656 20.996094 18.175781 21.296875 C 18.605469 21.664063 18.988281 22.398438 18.988281 23.515625 C 18.988281 25.121094 18.976563 26.414063 18.976563 26.808594 C 18.976563 27.128906 19.191406 27.503906 19.800781 27.386719 C 24.566406 25.796875 28 21.300781 28 16 C 28 9.371094 22.628906 4 16 4 Z"/>
            </svg>
        </footer>
        <p class="copyright">Copyright Branden Wheeler 2023</p>
    </div>
    <script src="scripts.js"></script>
</body>
</html>