import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Redirect } from 'react-router-dom';
import Main from './main';
import Home from './home';
import Explore from './explore';
import UserManage from './usermanage';
import BookDetail from './bookdetail';
import 'antd/dist/antd.css';
require('./css/main.css');

class AuthRoute extends React.Component {
    render() {
        const { component: Component, ...rest } = this.props;
        const isLogged = sessionStorage.getItem("userId") != null &&
            sessionStorage.getItem("userName") != null &&
            sessionStorage.getItem("userPermission") === this.props.permission &&
            sessionStorage.getItem("token") != null;

        if (!isLogged) {
            return <Route {...rest} render={(props) => <Redirect to='/' />} />
        } else {
            return <Route {...rest} render={this.props.render} />;
        }
    }
}

const StoreRouter = (
    <Router>
        <Route exact path='/' render={(props) => <Main {...props} content={Home} default={'0'} />} />
        <Route exact path='/explore' render={(props) => <Main {...props} content={Explore} default={'1'} />} />
        <AuthRoute permission='1' exact path='/usermanage' render={(props) => <Main {...props} content={UserManage} />} />
        <Route exact path='/book/:bookId' render={(props) => <Main {...props} content={BookDetail} />} />
    </Router>
);

ReactDOM.render(StoreRouter, document.getElementById('root'));
