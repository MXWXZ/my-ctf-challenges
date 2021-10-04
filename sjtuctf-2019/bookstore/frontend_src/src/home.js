import React, { Component } from 'react';
import axios from 'axios';

class Home extends Component {
    constructor(props) {
        super(props);

        this.state = {
            flag: "",
            userPermission: sessionStorage.getItem('userPermission'),
        };
    }

    componentDidMount() {
        if (this.state.userPermission === '1') {
            axios.get(`/api/flag`, { headers: { token: sessionStorage.getItem('token') } })
                .then(res => {
                    this.setState({
                        flag: res.data.data.flag,
                    })
                });
        }
    }

    render() {
        if (this.state.userPermission === '1') {
            return <h1 style={{ textAlign: 'center' }}>{this.state.flag}</h1>
        } else {
            return <h1 style={{ textAlign: 'center' }}>Where is the flag?</h1>
        }
    }
}

export default Home;
