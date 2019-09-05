import React from "react";
import styled from "@emotion/styled/macro";
import { NavLink } from "react-router-dom";

const Container = styled.aside`
  flex: 0 0 auto;
  position: sticky;
  height: 100vh;
  top: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-width: ${({ collapsed }) => (collapsed ? 0 : 30)}%;
  padding: 4rem ${({ collapsed }) => (collapsed ? 1 : 2)}rem;
  margin-right: -2rem;

  nav {
    overflow: hidden;
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

const SettingsBar = styled.div`
  display: flex;
  justify-content: space-evenly;
  overflow: hidden;

  a {
    display: flex;
    align-items: center;
    margin-right: 1rem;

    color: var(--secondary-color);
  }

  a:last-of-type {
    margin-right: 0;
  }

  a small {
    margin-left: 0.5rem;
  }
`;

const Sidepanel = ({ children, collapsed, onCollapse, onSignOut }) => {
  return (
    <Container collapsed={collapsed}>
      <CollapseButton onClick={onCollapse}>
        <i className="material-icons-round">
          {collapsed ? "arrow_right" : "arrow_left"}
        </i>
      </CollapseButton>

      {children}

      <SettingsBar>
        <NavLink to="/shares">
          <i className="material-icons-round">link</i>
          <small>Shares</small>
        </NavLink>
        <a href="#lougout" onClick={onSignOut}>
          <i className="material-icons-round">power_settings_new</i>
          <small>Logout</small>
        </a>
      </SettingsBar>
    </Container>
  );
};

export default Sidepanel;
