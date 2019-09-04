import React, { useState, useEffect } from "react";
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
import { ShareRoutes, GalleryRoutes } from "./Routes";
import SharesList from "./pages/shares";
import Sidepanel from "./components/Sidepanel";
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

const Content = styled.main`
  flex: 1;
  margin: 2rem;
  padding: 2rem 3rem;
  z-index: 1;

  background: white;
  border-radius: 8px;
  box-shadow: rgba(184, 194, 215, 0.25) 0px 4px 6px,
    rgba(184, 194, 215, 0.1) 0px 5px 7px;
`;

const App = ({ modules, authError, authToken, fetchModules, signOut }) => {
  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    apiClient.token = authToken;
    fetchModules();
  }, [authToken, fetchModules]);

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

  const renderSidebar = () => (
    <Sidepanel
      collapsed={collapsed}
      onCollapse={() => setCollapsed(!collapsed)}
      onSignOut={signOut}
    >
      <nav>{sidebarItems}</nav>
    </Sidepanel>
  );

  return (
    <BrowserRouter>
      <PageContainer>
        <Switch>
          <Route path="/gallery" render={renderSidebar} />
        </Switch>
        <Content>
          <Switch>
            <Route path="/share/:slug" component={ShareRoutes} />

            <Route path="/login" component={LoginForm} />
            {authError && <Redirect to="/login" />}

            {modules.length > 0 && (
              <Route path="/shares" component={SharesList} />
            )}

            {contentItems}
            {modules.length > 0 && <Redirect to={`/${modules[0]}`} />}

            <Route render={() => <h2>No active modules</h2>} />
          </Switch>
        </Content>
      </PageContainer>
    </BrowserRouter>
  );
};

export default connect(
  ({ modules, authError, authToken }) => ({ modules, authError, authToken }),
  { fetchModules, signOut }
)(App);
