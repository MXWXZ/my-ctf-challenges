import React, { Component } from 'react';
import { Tag, Row, Table } from 'antd';
import axios from 'axios';

class UserManage extends Component {
    state = {
        user: [],
    }

    componentDidMount() {
        axios.get(`/api/users`, { headers: { token: sessionStorage.getItem('token') } })
            .then(res => {
                this.setState({
                    user: res.data.data,
                })
            });
    }

    render() {
        const columns = [
            {
                title: 'ID',
                dataIndex: 'userId',
                align: 'center',
                sorter: (a, b) => a.userId - b.userId,
            }, {
                title: 'Username',
                dataIndex: 'userName',
                align: 'center',
                sorter: (a, b) => a.userName.localeCompare(b.userName),
            }, {
                title: 'Email',
                dataIndex: 'userEmail',
                align: 'center',
                sorter: (a, b) => a.userEmail.localeCompare(b.userEmail),
            }, {
                title: 'Permission',
                dataIndex: 'userPermission',
                align: 'center',
                sorter: (a, b) => a.userPermission - b.userPermission,
                render: userPermission => {
                    if (userPermission === 1)
                        return <Tag color='red'>admin</Tag>
                    else
                        return <Tag color='green'>user</Tag>
                }
            }, {
                title: 'Status',
                dataIndex: 'userStatus',
                align: 'center',
                sorter: (a, b) => a.userStatus - b.userStatus,
                render: userStatus => {
                    if (userStatus === 1)
                        return <Tag color='orange'>frozen</Tag>
                    else
                        return <Tag color='green'>normal</Tag>
                }
            }
        ];
        return (
            <Row>
                <Table rowKey={record => record.userId} bordered={true} columns={columns} dataSource={this.state.user} />
            </Row>
        );
    }
}

export default UserManage;
