
function ResetValidator() {
    this.modal = $('#set-password');
    this.modal.modal({ show: false, keyboard: false, backdrop: 'static' });
    this.alert = $('#set-password .alert');
    this.alert.hide();
}

function checkStrong(sValue) {
    var modes = 0;
    if (sValue.length < 1) return modes;
    if (/\d/.test(sValue)) modes++;
    if (/[a-z]/.test(sValue)) modes++;
    if (/[A-Z]/.test(sValue)) modes++;
    if (/\W/.test(sValue)) modes++;

    switch (modes) {
        case 1:
            return 1;
        case 2:
            return 2;
        case 3:
        case 4:
            return sValue.length < 6 ? 3 : 4
    }
}

ResetValidator.prototype.validatePassword = function (s) {
    if (checkStrong(s) == 4) {
        return true;
    } else {
        this.showAlert('Password needs digit, character(upper and lower case) and special characters and at least 6 length.');
        return false;
    }
}

ResetValidator.prototype.showAlert = function (m) {
    this.alert.attr('class', 'alert alert-danger');
    this.alert.html(m);
    this.alert.show();
}

ResetValidator.prototype.hideAlert = function () {
    this.alert.hide();
}

ResetValidator.prototype.showSuccess = function (m) {
    this.alert.attr('class', 'alert alert-success');
    this.alert.html(m);
    this.alert.fadeIn(500);
}