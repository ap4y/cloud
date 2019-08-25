import React, { Component } from "react";
import styled from "@emotion/styled";
import { connect } from "react-redux";
import {
  BrowserRouter,
  Route,
  NavLink,
  Switch,
  Redirect
} from "react-router-dom";

import AlbumsList from "./pages/albums";
import LoginForm from "./pages/login";
import { GalleryRoutes } from "./Routes";
import { apiClient, fetchModules, signOut } from "./actions";

const supportedModules = {
  gallery: {
    sidebar: AlbumsList,
    content: GalleryRoutes,
    icon: null
  }
};

const PageContainer = styled.div`
  min-height: 100vh;
  position: relative;

  display: flex;
  background: var(--primary-background-color);
`;

const Sidepanel = styled.aside`
  flex: 0 0 auto;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-width: ${({ collapsed }) => (collapsed ? 0 : 30)}%;
  height: 100vh;
  padding: 4rem 2rem;

  nav {
    position: sticky;
    top: 2rem;
    overflow: hidden;
  }

  div {
  }
`;

const SettingsBar = styled.div`
  display: flex;
  justify-content: center;
  overflow: hidden;

  a {
    display: flex;
    align-items: center;

    color: var(--secondary-color);
  }

  a small {
    margin-left: 0.5rem;
  }
`;

const CollapseButton = styled.a`
  display: flex;
  align-items: center;
  position: absolute;
  top: 3.5rem;
  right: -1rem;
  width: 3rem;
  height: 3rem;
  cursor: pointer;

  background: white;
  border-top-left-radius: 5px;
  border-bottom-left-radius: 5px;
  box-shadow: rgba(184, 194, 215, 0.25) 0px 4px 6px,
    rgba(184, 194, 215, 0.1) 0px 5px 7px;

  color: var(--secondary-color);
`;

const Content = styled.main`
  flex: 1;
  margin: 2rem 2rem 2rem 0;
  padding: 2rem 3rem;
  z-index: 1;

  background: white;
  border-radius: 8px;
  box-shadow: rgba(184, 194, 215, 0.25) 0px 4px 6px,
    rgba(184, 194, 215, 0.1) 0px 5px 7px;
`;

class App extends Component {
  state = { collapsed: false };

  componentDidMount() {
    apiClient.token = this.props.authToken;
    this.props.fetchModules();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.authToken === this.props.authToken) return;
    apiClient.token = this.props.authToken;
    this.props.fetchModules();
  }

  renderModules = () => {
    const { modules } = this.props;

    let navItems = [],
      sidebarItems = [],
      contentItems = [];
    modules.forEach(module => {
      const component = supportedModules[module];
      if (!component) return;

      navItems.push(
        <li key={module}>
          <NavLink to={module}>{component.icon}</NavLink>
        </li>
      );

      sidebarItems.push(
        <Route key={module} path={`/${module}`} component={component.sidebar} />
      );

      contentItems.push(
        <Route key={module} path={`/${module}`} component={component.content} />
      );
    });

    return { navItems, sidebarItems, contentItems };
  };

  render() {
    const { collapsed } = this.state;
    const { modules, authError } = this.props;
    const { navItems, sidebarItems, contentItems } = this.renderModules();

    return (
      <BrowserRouter>
        <PageContainer>
          <Sidepanel collapsed={collapsed || navItems.length === 0}>
            {navItems.length > 0 && (
              <CollapseButton
                onClick={() => this.setState({ collapsed: !collapsed })}
              >
                <i className="material-icons-round">
                  {collapsed ? "arrow_right" : "arrow_left"}
                </i>
              </CollapseButton>
            )}

            {navItems.length > 1 && (
              <nav>
                <ul>{navItems}</ul>
              </nav>
            )}

            <nav>{sidebarItems}</nav>

            <SettingsBar>
              <a href="#lougout" onClick={this.props.signOut}>
                <i className="material-icons-round">power_settings_new</i>{" "}
                <small>Logout</small>
              </a>
            </SettingsBar>
          </Sidepanel>

          <Content>
            <Switch>
              <Route path="/login" component={LoginForm} />
              {authError && <Redirect to="/login" />}
              {contentItems}
              {modules.length > 0 && <Redirect to={`/${modules[0]}`} />}
              <Route render={() => <h2>No active modules</h2>} />
            </Switch>
          </Content>
        </PageContainer>
      </BrowserRouter>
    );
  }
}

export default connect(
  ({ modules, authError, authToken }) => ({ modules, authError, authToken }),
  { fetchModules, signOut }
)(App);
