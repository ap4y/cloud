import React from "react";
import styled from "@emotion/styled/macro";

const CellItem = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  vertical-align: center;

  > i {
    font-size: ${({ large }) => (large ? 96 : 72)}px;
    text-align: center;
  }

  span {
    word-break: break-all;
    text-align: center;
  }
`;

export const DirCell = ({ dir: { name } }) => (
  <CellItem>
    <i className="material-icons-round">folder</i>
    <span>{name}</span>
  </CellItem>
);

export const FileCell = ({
  file: { name, updated_at },
  large,
  withModTime
}) => (
  <CellItem large={large}>
    <i className="material-icons-round">insert_drive_file</i>
    <span>{name}</span>
    {withModTime && (
      <span>
        <strong>Updated At: </strong>
        {new Date(updated_at).toLocaleString()}
      </span>
    )}
  </CellItem>
);
