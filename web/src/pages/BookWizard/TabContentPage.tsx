import React, { useState, useEffect } from 'react';
import { TabContent } from '../../models/Book';
import Editor from 'ckeditor5/build/ckeditor';
import { CKEditor } from "@ckeditor/ckeditor5-react";
import SimpleUploadAdapter from '@ckeditor/ckeditor5-upload/src/adapters/simpleuploadadapter';

type TabConentPageProps = {
  onContentChange: ( data: TabContent ) => void
};

export const TabContentPage: React.FC<TabConentPageProps> = ( {onContentChange}) => {

  const [video, setVideo] = useState< string | undefined>(undefined);
  const [body, setBody] = useState< string >("");
  useEffect( () => {
    onContentChange({
      video: video, 
      body: body,
    });
  }, [body, video]);

  const editorConfiguration = {
    toolbar: [ 'heading','|','bold','italic','underline','link','bulletedList','numberedList', '|','blockQuote','imageUpload','insertTable','undo','redo'],
    simpleUpload: {
      uploadUrl: "https://api.imgur.com/3/image",
    },
    image: {
      upload: {
        types: ['jpeg, png']
      }
    }
  };

  return (
    <div>
      <input 
        type="text"
        placeholder="Enter me"
        onChange={ e => setVideo(e.target.value)}
      />
      <CKEditor
        editor= { Editor }
        config = { editorConfiguration }
        data="<p>Hello from CKEditor 5!</p>"
        onChange= { (event: any, editor: any) => {
          const data = editor.getData();
          setBody(data);
        }}/>
    </div>
  );
};

