import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route } from 'react-router-dom';

import Header from './components/header'
import Footer from './components/footer'
import Home from './components/home'
import NewsOverview from './components/newsoverview'
import NewsEdit from './components/newsedit'
import FeatureSlides from './components/featureslides'
import KeukenDienst from './components/keukendienst'

const App = function () {
    return (
        <div>
        <Router>
            <div>
                <Header />
                <div className="container">
                    <Route path="/" exact component={Home} />
                    <Route path="/news" exact component={NewsOverview} />
                    <Route path="/news/edit/:id?" component={NewsEdit} />
                    <Route path="/featureslides" exact component={FeatureSlides}/>
                    <Route path="/keukendienst" exact component={KeukenDienst}/>
                </div>
                <Footer />
            </div>
        </Router>
    </div>
    );
};

ReactDOM.render(
    <App />,
    document.querySelector("#container")
);