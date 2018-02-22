import React, { Component } from 'react'
import { Row, Preloader, Col, Input, Button, Icon } from 'react-materialize'
import { withRouter } from "react-router-dom";
import { Editor } from 'react-draft-wysiwyg';
import { EditorState, convertToRaw, ContentState  } from 'draft-js'
import draftToHtml from 'draftjs-to-html';
import htmlToDraft from 'html-to-draftjs';
import DatePicker from 'react-datepicker'
import moment from "moment";
import 'react-draft-wysiwyg/dist/react-draft-wysiwyg.css'
import 'react-datepicker/dist/react-datepicker.css';

import { getNews, editNews, addNews } from '../apis/news_api'

const getToday = () => {
    let date = new Date()
    date.setUTCHours(0)
    date.setUTCMinutes(0)
    date.setUTCSeconds(0)
    date.setUTCMilliseconds(0)

    return moment(date)
}

class NewsEdit extends Component {
    constructor(props) {
        super(props)

        this.state = { title: "", id:-1, loading: false, editorState: EditorState.createEmpty(), isSaving: false, slideTime: 10, from: getToday(), to: getToday()}

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
        this.uploadImageCallBack = this.uploadImageCallBack.bind(this)
        this.setFromTime = this.setFromTime.bind(this)
        this.setToTime = this.setToTime.bind(this)
    }

    onGetNewsInfo(res) {
        const contentBlock = htmlToDraft(res.data.content);
        const contentState = ContentState.createFromBlockArray(contentBlock.contentBlocks);
        this.setState({ loading: false, title: res.data.title, editorState: EditorState.createWithContent(contentState), slideTime: res.data.slideTime, to: moment(new Date(res.data.to)), from: moment(new Date(res.data.from)) })
    }

    onEditorStateChange(editorState) {
        this.setState({
          editorState,
        });
    };

    save() {
        const content = draftToHtml(convertToRaw(this.state.editorState.getCurrentContent()))
        const title = this.title.input.value
        const slideTime = parseInt(this.slideTime.input.value, 10)

        if (!title) {
            window.Materialize.toast('Vul een titel in', 4000)
            return
        }

        this.state.id === -1 ? addNews({ title, content, slideTime, from: this.state.from, to: this.state.to }).then(this.doneSaving) : editNews({ id: this.state.id, title, content, slideTime, from: this.state.from, to: this.state.to }).then(this.doneSaving)
        this.setState({ isSaving: true })
    }

    doneSaving() {
        this.setState({ isSaving: false })
        this.props.history.push("/news");
    }

    setToTime(date) {
        this.setState({ to: date })
    }

    setFromTime(date) {
        console.log(date)
        this.setState({ from: date })
    }

    // credit to https://github.com/jpuri/react-draft-wysiwyg/blob/master/stories/ImageUpload/index.js
    uploadImageCallBack(file) {
        const HOST = `http://${window.location.host}`
        return new Promise(
          (resolve, reject) => {
            const xhr = new XMLHttpRequest(); // eslint-disable-line no-undef
            xhr.open('POST', `${HOST}/api/image`);
            const data = new FormData(); // eslint-disable-line no-undef
            data.append('image', file);
            data.append('host', HOST);
            xhr.send(data);
            xhr.addEventListener('load', () => {
              const response = JSON.parse(xhr.responseText);
              console.log(response)
              resolve({ data: response });
            });
            xhr.addEventListener('error', () => {
              const error = JSON.parse(xhr.responseText);
              reject(error);
            });
          },
        );
      }

    render() {
        if (this.state.loading) {
            return <Row><Col s={4} offset='s6'><Preloader size='big' flashing={true} /></Col></Row>
        }
        return <div className="marginBottom containerHeightNews">
            <Row><h2>Nieuws artikel bewerken</h2></Row>
            <Row>
	            <Input s={12} id="test" label="Titel" validate defaultValue={this.state.title} ref={(c) => this.title = c} />
            </Row>
            <Row>
                <Input type="number" label="Duur slide" s={4} defaultValue={this.state.slideTime} ref={(c) => this.slideTime = c} />
            </Row>
            <Row>
            <Editor
                toolbarClassName="toolbarClassName"
                wrapperClassName="wrapperClassName"
                editorClassName="editorClassName"
                editorState={this.state.editorState}
                onEditorStateChange={this.onEditorStateChange}
                toolbar={{
                    inline: { inDropdown: true },
                    list: { inDropdown: true },
                    textAlign: { inDropdown: true },
                    link: { inDropdown: true },
                    history: { inDropdown: true },
                    image: { uploadCallback: this.uploadImageCallBack, alt: { present: false, mandatory: false },
                    inputAccept: 'image/gif,image/jpeg,image/jpg,image/png,image/svg', },
                }}
            />
            </Row>

            <Row>
                <DatePicker label="Van" dateFormat="YYYY-MM-DD" selected={this.state.from} onChange={this.setFromTime}/>
            </Row>
            <Row>
                <DatePicker label="Tot" dateFormat="YYYY-MM-DD" selected={this.state.to} onChange={this.setToTime}/>
            </Row>
            
            <Row>
                <Col s={2}><Button waves='light' disabled={this.state.isSaving} onClick={this.save}>Opslaan<Icon left>save</Icon></Button></Col>
            </Row>
        </div>
    }
}

export default withRouter(NewsEdit)