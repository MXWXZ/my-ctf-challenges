
function AccountValidator() {
    // build array maps of the form inputs & control groups //

    this.formFields = [$('#name-tf'), $('#email-tf'), $('#user-tf'), $('#pass-tf')];
    this.controlGroups = [$('#name-cg'), $('#email-cg'), $('#user-cg'), $('#pass-cg')];

    // bind the form-error modal window to this controller to display any errors //

    this.alert = $('.modal-form-errors');
    this.alert.modal({ show: false, keyboard: true, backdrop: true });

    this.validateName = function (s) {
        return s.length >= 3;
    }

    this.validatePassword = function (s) {
        // if user is logged in and hasn't changed their password, return ok
        if ($('#userId').val() && s === '') {
            return true;
        } else {
            return s.length >= 6;
        }
    }

    this.validateEmail = function (e) {
        var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return re.test(e);
    }

    this.showErrors = function (a) {
        $('.modal-form-errors .modal-body p').text('Please correct the following problems :');
        var ul = $('.modal-form-errors .modal-body ul');
        ul.empty();
        for (var i = 0; i < a.length; i++) ul.append('<li>' + a[i] + '</li>');
        this.alert.modal('show');
    }

}

AccountValidator.prototype.showInvalidEmail = function () {
    this.controlGroups[1].addClass('error');
    this.showErrors(['That email address is already in use.']);
}

AccountValidator.prototype.showInvalidUserName = function () {
    this.controlGroups[2].addClass('error');
    this.showErrors(['That username is already in use.']);
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

AccountValidator.prototype.validateForm = function () {
    var e = [];
    for (var i = 0; i < this.controlGroups.length; i++) this.controlGroups[i].removeClass('error');
    if (this.validateName(this.formFields[0].val()) == false) {
        this.controlGroups[0].addClass('error'); e.push('Please Enter Your Name');
    }
    if (this.validateEmail(this.formFields[1].val()) == false) {
        this.controlGroups[1].addClass('error'); e.push('Please Enter A Valid Email');
    }
    if (this.validateName(this.formFields[2].val()) == false) {
        this.controlGroups[2].addClass('error');
        e.push('Please Choose A Username');
    }
    if (this.validatePassword(this.formFields[3].val()) == false || checkStrong(this.formFields[3].val()) != 4) {
        this.controlGroups[3].addClass('error');
        e.push('Password needs digit, character(upper and lower case) and special characters and at least 6 length.');
    }
    if (e.length) this.showErrors(e);
    return e.length === 0;
}

