import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon } from 'react-materialize'
import { withRouter } from "react-router-dom";
import { Editor } from 'react-draft-wysiwyg'
import { EditorState, convertToRaw, ContentState  } from 'draft-js'
import draftToHtml from 'draftjs-to-html';
import htmlToDraft from 'html-to-draftjs';
import 'react-draft-wysiwyg/dist/react-draft-wysiwyg.css'

import { getNews, editNews, addNews } from '../apis/news_api'

class NewsEdit extends Component {
    constructor(props) {
        super(props)
        this.state = { title: "", id:-1, loading: false, editorState: EditorState.createEmpty(), isSaving: false }

        if (this.props.match.params.id !== "new") {
            this.state.loading = true
            this.state.id = parseInt(this.props.match.params.id, 10)
            getNews(this.props.match.params.id).then(this.onGetNewsInfo.bind(this))          
        } else {
            const contentBlock = htmlToDraft("");
            const contentState = ContentState.createFromBlockArray(contentBlock.contentBlocks);
            this.state.editorState = EditorState.createWithContent(contentState); 
        }

        this.onEditorStateChange = this.onEditorStateChange.bind(this)
        this.save = this.save.bind(this)
        this.doneSaving = this.doneSaving.bind(this)
    }

    onGetNewsInfo(res) {
        const contentBlock = htmlToDraft(res.data.content);
        const contentState = ContentState.createFromBlockArray(contentBlock.contentBlocks);

        this.setState({ loading: false, title: res.data.title, editorState: EditorState.createWithContent(contentState) })
    }

    onEditorStateChange(editorState) {
        this.setState({
          editorState,
        });
    };

    save() {
        const content = draftToHtml(convertToRaw(this.state.editorState.getCurrentContent()))
        console.log(this.title.input);
        const title = this.title.input.value

        this.state.id === -1 ? addNews({ title, content }).then(this.doneSaving) : editNews({ id: this.state.id, title, content }).then(this.doneSaving)
        this.setState({ isSaving: true })
    }

    doneSaving() {
        this.setState({ isSaving: false })
        this.props.history.push("/news");
    }

    render() {
        if (this.state.loading) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }
        return <div>
            <Row><h2>Nieuws artikel bewerken</h2></Row>
            <Row>
	            <Input s={12} label="Titel" validate defaultValue={this.state.title} ref={(c) => this.title = c} />
            </Row>
            <Row>
            <Editor
                toolbarClassName="toolbarClassName"
                wrapperClassName="wrapperClassName"
                editorClassName="editorClassName"
                editorState={this.state.editorState}
                onEditorStateChange={this.onEditorStateChange}
            />
            </Row>
            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default withRouter(NewsEdit)