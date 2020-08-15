import React from "react";
import styled from "@emotion/styled/macro";

const FileList = styled.ul`
  list-style: none;

  li {
    display: flex;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  i {
    margin-right: 0.5rem;
  }
`;

const FilesItems = ({ items }) => {
  return (
    <FileList>
      {items.map((item, idx) => (
        <li key={idx}>
          <i className="material-icons-round">arrow_right</i>
          {item}
        </li>
      ))}
    </FileList>
  );
};

export default FilesItems;
