import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Dropdown, Button, Menu, Icon } from 'antd';

function handleMenuClick(e) {
    if (e.key === '0') {
        sessionStorage.clear();
        window.location.reload();
    }
}

const userMenu = (
    <Menu onClick={handleMenuClick}>
        <Menu.Item key='0' >Logout</Menu.Item>
    </Menu>
);

const adminMenu = (
    <Menu onClick={handleMenuClick}>
        <Menu.Item key='1'><Link to='/usermanage'>User Manage</Link></Menu.Item>
        <Menu.Item key='0'>Logout</Menu.Item>
    </Menu>
);

/*
    AvatarBar
    @param  userName: username
    @param  userPermission: 0 for user 1 for admin
*/
class AvatarBar extends Component {
    render() {
        return (
            <div>
                <Dropdown overlay={this.props.userPermission === '0' ? userMenu : adminMenu} trigger={['click']}>
                    <Button className='nav-button' icon='user'>
                        {this.props.userName}
                        <Icon type='down' />
                    </Button>
                </Dropdown>
                <Link to='/'>
                    <Button className='nav-button' type="primary" shape="circle" icon="shopping-cart"
                        style={{ display: this.props.userPermission === '1' ? 'none' : 'inline-block' }} />
                </Link>
            </div>
        );
    }
}

export default AvatarBar;
