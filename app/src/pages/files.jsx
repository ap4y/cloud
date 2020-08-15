import React, { useState, useEffect, useRef } from "react";
import styled from "@emotion/styled/macro";
import { connect } from "react-redux";
import { NavLink } from "react-router-dom";

import { DirCell, FileCell } from "../components/files/Cells";
import FilesToolbar from "../components/files/Toolbar";
import {
  apiClient,
  fetchFile,
  removeFile,
  uploadFile,
  createFolder,
  removeFolder,
  shareFolder
} from "../lib/actions";
import { locateInTree } from "../lib/utils";
import useSharing from "../lib/useSharing";

const viewableExt = ["md", "org", "txt"];

const Files = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  grid-gap: 2rem;
`;

const FilesItem = styled.div`
position: relative;

  a {
    color: var(--secondary-color);
  }

  a:hover {
    color: var(--primary-color);
  }

  &:before {
    display: ${({ active }) => (active ? "block" : "none")};
    content: "";
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
    cursor: pointer;
  }

  &:after {
    display: ${({ active }) => (active ? "flex" : "none")};
    content: "${({ selected }) => (selected ? "✓" : "")}";
    position: absolute;
    top: 1rem;
    right: 1rem;
    width: 25px;
    height: 25px;
    justify-content: center;
    align-items: center;
    z-index: 2;

    background: var(--outline-color);
    border-radius: 15px;

    color: var(--primary-text-color);
    font-weight: 700;
  }
`;

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
  share,
  match,
  history,
  fetchFile,
  removeFile,
  uploadFile,
  createFolder,
  removeFolder,
  shareFolder
}) => {
  const [content, setContent] = useState(null);
  const [
    sharePopup,
    showSharing,
    setShowSharing,
    sharedItems,
    toggleSharedItem
  ] = useSharing(folder.children, expireAt =>
    shareFolder(folder, sharedItems, expireAt)
  );

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

    fetchFile(file.url, share).then(content => setContent(content));
  }, [file, fetchFile, share]);

  useEffect(() => {
    setShowSharing(false);
  }, [folder, setShowSharing]);

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

  const mkdir = e => {
    e.preventDefault();
    const folderName = window.prompt("Folder Name");
    if (!folderName) return;

    createFolder(folder, folderName);
  };

  const rmdir = e => {
    e.preventDefault();
    if (!window.confirm("Are you sure?")) return;

    const parentPath = folder.path
      .split("/")
      .slice(0, -1)
      .join("/");
    removeFolder(folder).then(() => history.replace(`/files${parentPath}`));
  };

  const toggleSharing = e => {
    e.preventDefault();
    setShowSharing(true);
  };

  const fileItems = folder.children
    .sort((a, b) => `${a.type}${a.name}`.localeCompare(`${b.type}${b.name}`))
    .map(file => {
      return (
        <FilesItem
          key={file.path}
          active={showSharing}
          selected={sharedItems.includes(file)}
          onClick={() => toggleSharedItem(file)}
        >
          <NavLink to={`${match.url}/${encodeURI(file.name)}`}>
            {file.type === "directory" && <DirCell dir={file} />}
            {file.type === "file" && <FileCell file={file} />}
          </NavLink>
        </FilesItem>
      );
    });

  return (
    <FilesContent>
      <FilesToolbar
        path={file ? file.path : folder.path}
        file={file}
        fileURL={file && apiClient.fileURL(file, share)}
        canSave={content !== null}
        canEdit={share === undefined}
        onSave={saveFile}
        onDelete={deleteFile}
        onUpload={presentUpload}
        onMkdir={mkdir}
        onRmdir={rmdir}
        onShare={toggleSharing}
      />
      {showSharing && sharePopup}
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
    const { path = "", slug } = props.match.params;
    const { folder, file } = locateInTree(tree, path);
    return { file, folder, items: folder.children, share: slug };
  },
  { fetchFile, removeFile, uploadFile, createFolder, removeFolder, shareFolder }
)(FilesGrid);
