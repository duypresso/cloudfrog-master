import React, { useRef, useState } from 'react';
import { CloudArrowUpIcon } from '@heroicons/react/24/outline';
import { motion } from 'framer-motion';

const FileUploader = ({ file, setFile, onDrop }) => {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef(null);

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      setFile(selectedFile);
    }
  };

  const handleDragOver = (e) => {
    e.preventDefault();
  };

  const handleDrop = (e) => {
    e.preventDefault();
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      setFile(droppedFile);
      if (onDrop) onDrop(droppedFile);
    }
  };

  return (
    <motion.div 
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
      className={`border-2 border-dashed rounded-lg p-8 mb-6 text-center transition-all duration-300
        ${isDragging ? 'border-blue-400 bg-blue-50' : ''}
        ${file ? 'border-green-300 bg-green-50' : 'border-gray-300'}`}
      onDragEnter={() => setIsDragging(true)}
      onDragLeave={() => setIsDragging(false)}
      onDrop={(e) => {
        setIsDragging(false);
        handleDrop(e);
      }}
      onDragOver={handleDragOver}
    >
      <CloudArrowUpIcon className={`h-16 w-16 mx-auto ${file ? 'text-green-500' : 'text-gray-400'} transition-colors duration-300`} />
      
      <p className="mt-4 text-sm text-gray-600">
        Drag and drop your file here, or
      </p>
      
      <button
        onClick={() => fileInputRef.current.click()}
        className="mt-2 text-blue-500 hover:text-blue-700 font-medium"
      >
        browse to upload
      </button>
      
      <input
        ref={fileInputRef}
        type="file"
        className="hidden"
        onChange={handleFileChange}
      />

      {file && (
        <div className="mt-4 p-3 bg-white rounded-lg shadow-sm border border-green-100 animate-slideUp">
          <div className="flex items-center">
            <div className="bg-green-100 rounded-full p-2 mr-3">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-green-600" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="flex-1 truncate">
              <p className="text-sm font-medium text-gray-700 truncate">{file.name}</p>
              <p className="text-xs text-gray-500">
                {(file.size / 1024 / 1024).toFixed(2)} MB
              </p>
            </div>
            <button
              onClick={() => setFile(null)}
              className="text-gray-400 hover:text-gray-600 ml-2"
              title="Remove file"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 111.414 1.414L11.414 10l4.293 4.293a1 1 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 01-1.414-1.414L8.586 10 4.293 5.707a1 1 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      )}
    </motion.div>
  );
};

export default FileUploader;
