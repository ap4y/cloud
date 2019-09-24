import React from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

const DirItem = styled.li`
  padding-left: ${({ level }) => level}rem;

  a {
    display: flex;
    align-items: center;
  }

  a span {
    margin-left: 1rem;

    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
`;

export const DirTree = ({ leaf: { path, name, children }, level }) => (
  <>
    <DirItem level={level}>
      <NavLink to={`/files${path}/`} exact>
        <i className="material-icons-round">folder</i>
        <span>{name}</span>
      </NavLink>
    </DirItem>
    {children
      .filter(({ type }) => type === "directory")
      .sort((a, b) => (a.name > b.name ? 1 : -1))
      .map(item => (
        <DirTree key={item.path} leaf={item} level={level + 1} />
      ))}
  </>
);

const Tree = styled.ul`
  display: flex;
  flex-direction: column;
  list-style: none;
`;

export const FilesTree = ({ tree }) => {
  return (
    <div>
      <h2>Files</h2>
      <Tree>
        <DirTree leaf={tree} level={0} />
      </Tree>
    </div>
  );
};

export default connect(
  ({ files: { tree } }) => ({ tree }),
  {}
)(FilesTree);
