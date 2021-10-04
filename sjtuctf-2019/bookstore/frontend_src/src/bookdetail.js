import React, { Component } from 'react';
import { withRouter } from 'react-router';
import { message, Typography, Row, Col, InputNumber, Button } from 'antd';
import axios from 'axios';

const { Text } = Typography;

class BookDetail extends Component {
    state = {
        bookId: 0,
        bookName: '',
        bookAuthor: '',
        bookPrice: 0,
        bookIsbn: '',
        bookAmount: 1,
        bookDetail: '',
        buyAmount: 1,
    }

    constructor(props) {
        super(props);
        this.addCart = this.addCart.bind(this);
        this.onAmountChange = this.onAmountChange.bind(this);
    }

    componentDidMount() {
        const urlParam = this.props.match.params;
        axios.get(`/api/books`, {
            params: {
                bookId: urlParam.bookId
            }
        })
            .then(res => {
                const data = res.data.data;
                this.setState({
                    bookId: data.bookId,
                    bookName: data.bookName,
                    bookAuthor: data.bookAuthor,
                    bookPrice: data.bookPrice,
                    bookIsbn: data.bookIsbn,
                    bookAmount: data.bookAmount,
                    bookDetail: data.bookDetail
                })
            });
    }

    addCart() {
        if (!sessionStorage.getItem('userId')) {
            document.getElementById('signin').click();
        } else {
            message.success("Add to cart successfully!");
        }
    }

    onAmountChange(value) {
        let num = parseInt(value);
        this.setState({ buyAmount: isNaN(num) ? 1 : num })
    }

    render() {
        let operation;
        if (sessionStorage.getItem('userPermission') !== '1')
            operation = (
                <div style={{ display: 'inline' }}>
                    <Row style={{ marginTop: '10px' }}>
                        <Text style={{ marginRight: '10px' }}>Amount: </Text>
                        <InputNumber id='amount' min={1} max={this.state.bookAmount === 0 ? 1 : this.state.bookAmount} defaultValue={1} onChange={this.onAmountChange} />
                        <Text style={{ marginLeft: '20px' }}>Stock: {this.state.bookAmount}</Text>
                    </Row>
                    <Row style={{ marginTop: '20px' }}>
                        <Button size='large' onClick={this.addCart} disabled={this.state.bookAmount === 0 ? true : false}>Add to cart</Button>
                    </Row>
                </div>
            );

        return (
            <Row>
                <Col span={6} className='book-detail-img' style={{ height: '350px' }}>
                    {
                        <img alt='cover' src='/img/he.jpg' style={{ height: '350px', width: '100%' }} />
                    }
                </Col>
                <Col span={18} className='book-detail-info'>
                    <div style={{ height: '180px' }}>
                        <Row>
                            <h1 style={{ float: 'left' }}>{this.state.bookName}</h1>
                            <Text type='secondary' style={{ float: 'right' }}>ISBN: {this.state.bookIsbn}</Text>
                            <div className='clear' />
                        </Row>
                        <Row>
                            <Text type='secondary'>{this.state.bookAuthor}</Text>
                        </Row>
                        <Row style={{ marginTop: '50px' }}>
                            <Text>{this.state.bookDetail}</Text>
                        </Row>
                    </div>
                    <div style={{ height: '170px' }}>
                        <Row style={{ paddingTop: '10px' }}>
                            <Text className='book-price'>&yen; {(this.state.bookPrice / 100 * this.state.buyAmount).toFixed(2)}</Text>
                        </Row>
                        {operation}
                    </div>
                </Col>
            </Row>
        );
    }
}

export default withRouter(BookDetail);