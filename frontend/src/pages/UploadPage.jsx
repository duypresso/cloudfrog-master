import React, { useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { ClipboardIcon } from '@heroicons/react/24/outline';
import FileUploader from '../components/FileUploader';
import ProgressBar from '../components/ProgressBar';
import { showToast } from '../utils/toast';
import api from '../utils/api';

const MAX_FILE_SIZE = 100 * 1024 * 1024; // 100MB

const UploadPage = () => {
  const [file, setFile] = useState(null);
  const [progress, setProgress] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [downloadUrl, setDownloadUrl] = useState('');
  const [error, setError] = useState('');

  const copyToClipboard = () => {
    navigator.clipboard.writeText(downloadUrl)
      .then(() => showToast.success('Link copied to clipboard!'));
  };

  const validateFile = (file) => {
    if (!file) return 'Please select a file.';
    if (file.size > MAX_FILE_SIZE) {
      return 'File size exceeds 100MB limit.';
    }
    return null;
  };

  const handleUpload = async () => {
    const errorMessage = validateFile(file);
    if (errorMessage) {
      setError(errorMessage);
      showToast.error(errorMessage);
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    const toastId = showToast.loading('Uploading file...');
    setUploading(true);
    setProgress(0);
    setDownloadUrl('');
    setError('');

    try {
      const response = await api.post('/upload', formData, {
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          setProgress(percentCompleted);
        },
      });

      setDownloadUrl(response.data.downloadUrl);
      showToast.success('File uploaded successfully!');
    } catch (err) {
      const errorMessage = err.message || 'Upload failed. Please try again.';
      setError(errorMessage);
      showToast.error(errorMessage);
      console.error('Upload error details:', err.details);
    } finally {
      setUploading(false);
      toast.dismiss(toastId);
    }
  };

  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-2xl mx-auto"
    >
      <div className="card mb-8">
        <h1 className="text-3xl font-bold text-center mb-6 text-gray-800">
          Share Files Securely
        </h1>
        
        <p className="text-gray-600 text-center mb-8">
          Upload files up to 100MB and share them with anyone using a secure link that automatically expires.
        </p>

        <FileUploader file={file} setFile={setFile} />

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6 animate-fadeIn">
            <div className="flex">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-red-500" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
              {error}
            </div>
          </div>
        )}

        {uploading && <ProgressBar progress={progress} />}

        {downloadUrl && (
          <div className="bg-green-50 border border-green-100 rounded-lg p-5 mb-6 animate-slideUp">
            <div className="flex items-center text-green-700 mb-3">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <h3 className="text-lg font-medium">File Uploaded Successfully!</h3>
            </div>
            
            <p className="text-sm text-gray-600 mb-3">
              Your file is now available at this link:
            </p>
            
            <div className="flex items-center">
              <input
                type="text"
                value={downloadUrl}
                readOnly
                className="flex-grow p-2 border rounded-l-lg outline-none bg-white text-gray-700"
              />
              <button
                onClick={copyToClipboard}
                className="bg-blue-500 text-white p-2 rounded-r-lg hover:bg-blue-600 transition-colors"
                title="Copy to clipboard"
              >
                <ClipboardIcon className="h-5 w-5" />
              </button>
            </div>
            
            <p className="text-sm text-gray-500 mt-3 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
              </svg>
              This link will expire in 7 days
            </p>
          </div>
        )}

        <button
          onClick={handleUpload}
          disabled={!file || uploading}
          className={`w-full py-3 px-4 rounded-lg font-medium transition-all duration-200
            ${
              !file || uploading
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'bg-blue-500 text-white hover:bg-blue-600 shadow-sm hover:shadow'
            }`}
        >
          {uploading ? 'Uploading...' : 'Upload File'}
        </button>
      </div>
      
      <div className="card">
        <h2 className="text-xl font-semibold mb-3 text-gray-800">About CloudFrog</h2>
        <p className="text-gray-600 mb-4">
          CloudFrog is a secure way to share files with others. All files are stored with encryption and links automatically expire to keep your data safe.
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="p-4 border border-gray-100 rounded-lg bg-gray-50">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-blue-500 mb-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clipRule="evenodd" />
            </svg>
            <h3 className="font-medium mb-1">Easy Sharing</h3>
            <p className="text-sm text-gray-600">Share files with anyone using a simple link</p>
          </div>
          
          <div className="p-4 border border-gray-100 rounded-lg bg-gray-50">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-blue-500 mb-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
            </svg>
            <h3 className="font-medium mb-1">Secure Storage</h3>
            <p className="text-sm text-gray-600">Files are stored securely and access is controlled</p>
          </div>
          
          <div className="p-4 border border-gray-100 rounded-lg bg-gray-50">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-blue-500 mb-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
            </svg>
            <h3 className="font-medium mb-1">Auto Expiry</h3>
            <p className="text-sm text-gray-600">Links expire automatically for better security</p>
          </div>
        </div>
      </div>
    </motion.div>
  );
};

export default UploadPage;
