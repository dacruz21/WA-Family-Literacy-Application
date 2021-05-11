import React, { useState, useEffect } from 'react';
import { Book, CreateBook, TabContent } from '../models/Book';
import { TabContentPage } from './BookWizard/TabContentPage';

export const BookWizardPage: React.FC = () => {

  const emptyTabContent: TabContent = {
    body: "",
    video: undefined
  };

  // const [readTabContent, setReadTabContent] =  useState<TabContent | null>(null);
  const [title, setTitle] = useState<string>("");
  const [author, setAuthor] = useState<string>("");
  const [image, setImage] = useState<Uint8Array>(new Uint8Array());
  const [readTabContent, setReadTabContent] = useState<TabContent>(emptyTabContent);
  const [exploreTabContent, setExploreTabContent] = useState<TabContent>(emptyTabContent);
  const [learnTabContent, setLearnTabContent] = useState<TabContent>(emptyTabContent);

  const updateReadTabContent = (data: TabContent): void => {
    setReadTabContent(data);
  };
  


  return (
    <div>
      <TabContentPage onContentChange= {updateReadTabContent}>
      </TabContentPage>
    </div>
  );

};