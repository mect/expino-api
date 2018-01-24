import React, { Component } from 'react';
import { Row, Collection, CollectionItem, Button , Icon} from 'react-materialize';
import { Link } from 'react-router-dom'
import { getAllNews, deleteNews } from '../apis/news_api' 
import _ from 'lodash/collection'

class NewsEdit extends Component {
    constructor(props) {
        super(props)
        this.state = { items: [] }

        getAllNews().then(this.gotNews.bind(this))
        this.deleteItem = this.deleteItem.bind(this)
    }

    gotNews(result) {
        this.setState({ items: result.data })
    }

    deleteItem(id) {
        this.setState({ items: _.reject(this.state.items, (i) => i.id === id) })
        deleteNews(id)
    }

    render() {
        const items = this.state.items.map(i => <CollectionItem>
                <Link to={`/news/edit/${i.id}`}>{i.title}</Link>
                <a onClick={() => this.deleteItem(i.id)}><Icon right>delete</Icon></a>
            </CollectionItem>)
        return <Row>
        <Collection header='Nieuws items'>
	        {items}
        </Collection>
        <Link to="/news/edit/new"><Button floating fab='vertical' icon='add' className='red' large style={{bottom: '45px', right: '24px'}}/></Link>
    </Row>
    }
}

export default NewsEdit