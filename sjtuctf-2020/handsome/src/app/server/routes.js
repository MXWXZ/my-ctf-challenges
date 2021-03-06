
var CT = require('./modules/country-list');
var AM = require('./modules/account-manager');
var EM = require('./modules/email-dispatcher');

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


module.exports = function (app) {

    /*
        login & logout
    */

    app.get('/', function (req, res) {
        // check if the user has an auto login key saved in a cookie //
        if (req.cookies.login == undefined) {
            res.render('login', { title: 'Hello - Please Login To Your Account' });
        } else {
            // attempt automatic login //
            AM.validateLoginKey(req.cookies.login, req.ip, function (e, o) {
                if (o) {
                    AM.autoLogin(o.user, o.pass, function (o) {
                        req.session.user = o;
                        res.redirect('/home');
                    });
                } else {
                    res.render('login', { title: 'Hello - Please Login To Your Account' });
                }
            });
        }
    });

    app.post('/', function (req, res) {
        AM.manualLogin(req.body['user'], req.body['pass'], function (e, o) {
            if (!o) {
                res.status(400).send(e);
            } else {
                req.session.user = o;
                if (req.body['remember-me'] == 'false') {
                    res.status(200).send(o);
                } else {
                    AM.generateLoginKey(o.user, req.ip, function (key) {
                        res.cookie('login', key, { maxAge: 900000 });
                        res.status(200).send(o);
                    });
                }
            }
        });
    });

    app.post('/logout', function (req, res) {
        res.clearCookie('login');
        req.session.destroy(function (e) { res.status(200).send('ok'); });
    })

    /*
        control panel
    */

    app.get('/home', function (req, res) {
        if (req.session.user == null) {
            res.redirect('/');
        } else if (req.session.user.email == 'admin@0ops.sjtu.cn') {
            res.render('admin', {
                title: 'Control Panel',
                udata: req.session.user
            });
        } else {
            res.render('home', {
                title: 'Control Panel',
                countries: CT,
                udata: req.session.user
            });
        }
    });

    app.post('/home', function (req, res) {
        if (req.session.user == null || req.session.user.email != 'admin@0ops.sjtu.cn') {
            res.redirect('/');
        } else {
            const pug = require('pug');
            res.status(200).send(pug.render(req.body['data']));
        }
    });

    // app.post('/home', function (req, res) {
    //     if (req.session.user == null) {
    //         res.redirect('/');
    //     } else {
    //         AM.updateAccount({
    //             id: req.session.user._id,
    //             name: req.body['name'],
    //             email: req.body['email'],
    //             pass: req.body['pass'],
    //             country: req.body['country']
    //         }, function (e, o) {
    //             if (e) {
    //                 res.status(400).send('error-updating-account');
    //             } else {
    //                 req.session.user = o.value;
    //                 res.status(200).send('ok');
    //             }
    //         });
    //     }
    // });

    /*
        new accounts
    */

    app.get('/signup', function (req, res) {
        res.render('signup', { title: 'Signup', countries: CT });
    });

    app.post('/signup', function (req, res) {
        if (checkStrong(req.body['pass']) != 4) {
            res.status(400).send('password too weak');
        } else {
            AM.addNewAccount({
                name: req.body['name'],
                email: req.body['email'],
                user: req.body['user'],
                pass: req.body['pass'],
                country: req.body['country']
            }, function (e) {
                if (e) {
                    res.status(400).send(e);
                } else {
                    res.status(200).send('ok');
                }
            });
        }
    });

    /*
        password reset
    */

    app.post('/lost-password', function (req, res) {
        let email = req.body['email'];
        AM.generatePasswordKey(email, req.ip, function (e, account) {
            if (e) {
                res.status(400).send(e);
            } else {
                EM.dispatchResetPasswordLink(account, function (l) {
                    // TODO this callback takes a moment to return, add a loader to give user feedback //
                    if (!e) {
                        res.status(200).send(l);
                    } else {
                        res.status(400).send('unable to dispatch password reset');
                    }
                });
            }
        });
    });

    app.get('/reset-password', function (req, res) {
        AM.validatePasswordKey(req.query['key'], req.ip, function (e, o) {
            if (e || o == null) {
                res.redirect('/');
            } else {
                req.session.passKey = req.query['key'];
                res.render('reset', { title: 'Reset Password' });
            }
        })
    });

    app.post('/reset-password', function (req, res) {
        if (checkStrong(req.body['pass']) != 4) {
            res.status(400).send('password too weak');
        } else {
            let newPass = req.body['pass'];
            let passKey = req.session.passKey;
            // destory the session immediately after retrieving the stored passkey //
            req.session.destroy();
            AM.updatePassword(passKey, newPass, function (e, o) {
                if (o) {
                    res.status(200).send('ok');
                } else {
                    res.status(400).send('unable to update password');
                }
            })
        }
    });

    /*
        view, delete & reset accounts
    */

    app.get('/whoami', function (req, res) {
        res.render('handsome');
    });

    app.get('/admin', function (req, res) {
        AM.getAllRecords(function (e, accounts) {
            res.render('print', { title: 'Account List', accts: accounts });
        })
    });

    app.get('*', function (req, res) { res.render('404', { title: 'Page Not Found' }); });

};
