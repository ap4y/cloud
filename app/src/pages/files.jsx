import React, { useState, useEffect, useRef } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

import { Toolbar } from "../components/Controls";
import { apiClient, fetchFile, removeFile, uploadFile } from "../actions";
import { locateInTree } from "../lib/utils";

const viewableExt = ["md", "org", "txt"];

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
  file: { name, path, updated_at },
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

const Files = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  grid-gap: 2rem;
`;

const FilesItem = styled.div`
  a {
    color: var(--secondary-color);
  }

  a:hover {
    color: var(--primary-color);
  }
`;

const FilesToolbar = ({
  path,
  file,
  canSave,
  onSave,
  onDelete,
  onUpload,
  onShare
}) => {
  return (
    <Toolbar title={file ? file.path : path}>
      {canSave && (
        <a href="#save" onClick={onSave}>
          <i className="material-icons-round">save</i>
        </a>
      )}
      {file && (
        <a href={apiClient.fileURL(file)} download={file.name}>
          <i className="material-icons-round">download</i>
        </a>
      )}
      {file && (
        <a href="#delete" onClick={onDelete}>
          <i className="material-icons-round">delete</i>
        </a>
      )}
      {!file && (
        <a href="#upload" onClick={onUpload}>
          <i className="material-icons-round">upload</i>
        </a>
      )}
      {!file && (
        <a href="#share" onClick={onShare}>
          <i className="material-icons-round">share</i>
        </a>
      )}
    </Toolbar>
  );
};

const TextEditor = styled.textarea`
  flex: 1 1 auto;
`;

const FilesContent = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;
`;

export const FilesGrid = ({
  file,
  folder,
  history,
  fetchFile,
  removeFile,
  uploadFile
}) => {
  const [content, setContent] = useState(null);

  const fileInputRef = useRef(null);

  useEffect(() => {
    if (!file) {
      setContent(null);
      return;
    }

    const [, ext] = file.name.split(".");
    if (!viewableExt.includes(ext)) {
      setContent(null);
      return;
    }

    fetchFile(file.url).then(content => setContent(content));
  }, [file, fetchFile]);

  const deleteFile = e => {
    e.preventDefault();
    if (!window.confirm("Are you sure?")) return;
    removeFile(folder, file).then(() =>
      history.replace(`/files${folder.path}`)
    );
  };

  const presentUpload = e => {
    e.preventDefault();
    fileInputRef.current.click();
  };

  const performUpload = ({ target: { files } }) => {
    if (files.length === 0) return;
    uploadFile(folder, files[0]);
  };

  const saveFile = e => {
    e.preventDefault();
    uploadFile(folder, new File([content], file.name));
  };

  const fileItems = folder.children
    .sort((a, b) => `${a.type}${a.name}`.localeCompare(`${b.type}${b.name}`))
    .map(file => {
      return (
        <FilesItem key={file.path}>
          {file.type === "directory" && (
            <NavLink to={`/files${file.path}/`}>
              <DirCell dir={file} />
            </NavLink>
          )}
          {file.type === "file" && (
            <NavLink to={`/files${file.path}`}>
              <FileCell file={file} />
            </NavLink>
          )}
        </FilesItem>
      );
    });

  return (
    <FilesContent>
      <FilesToolbar
        path={folder.path}
        file={file}
        canSave={file && content}
        onSave={saveFile}
        onDelete={deleteFile}
        onUpload={presentUpload}
      />
      {!file && <Files>{fileItems}</Files>}
      {file && !content && <FileCell file={file} large withModTime />}
      {file && content && (
        <TextEditor
          value={content}
          onChange={({ target }) => setContent(target.value)}
        />
      )}
      <input
        type="file"
        ref={fileInputRef}
        style={{ display: "none" }}
        onChange={performUpload}
      />
    </FilesContent>
  );
};

export default connect(
  ({ files: { tree } }, props) => {
    const { path = "" } = props.match.params;
    const { folder, file } = locateInTree(tree, path);
    return { file, folder, items: folder.children };
  },
  { fetchFile, removeFile, uploadFile }
)(FilesGrid);
