import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import api from '../utils/api';
import { ArrowPathIcon, ExclamationTriangleIcon, ArrowTopRightOnSquareIcon } from '@heroicons/react/24/outline';

const DownloadPage = () => {
  const { shortCode } = useParams();
  const [status, setStatus] = useState('loading'); // loading, error, expired, success

  useEffect(() => {
    const downloadFile = async () => {
      try {
        // In a real app, we might want to check if the file exists first
        window.location.href = `${api.defaults.baseURL}/download/${shortCode}`;
        setStatus('success');
      } catch (err) {
        if (err.response && err.response.status === 410) {
          setStatus('expired');
        } else {
          setStatus('error');
        }
      }
    };

    // Start download after a short delay
    const timer = setTimeout(() => {
      downloadFile();
    }, 1000);

    return () => clearTimeout(timer);
  }, [shortCode]);

  if (status === 'loading') {
    return (
      <div className="max-w-lg mx-auto text-center py-16 animate-fadeIn">
        <div className="card">
          <div className="rounded-full bg-blue-50 p-4 w-16 h-16 flex items-center justify-center mx-auto mb-6">
            <ArrowPathIcon className="h-8 w-8 text-blue-500 animate-spin" />
          </div>
          
          <h2 className="text-2xl font-bold mb-3 text-gray-800">Preparing Your Download</h2>
          
          <div className="h-1.5 w-full bg-gray-200 rounded-full mb-4 overflow-hidden">
            <div className="h-full bg-blue-500 rounded-full animate-pulse" style={{ width: '75%' }}></div>
          </div>
          
          <p className="text-gray-600">
            Your download will begin automatically in a moment...
          </p>
        </div>
      </div>
    );
  }

  if (status === 'error') {
    return (
      <div className="max-w-lg mx-auto text-center py-16 animate-fadeIn">
        <div className="card border-red-100">
          <div className="rounded-full bg-red-50 p-4 w-16 h-16 flex items-center justify-center mx-auto mb-6">
            <ExclamationTriangleIcon className="h-8 w-8 text-red-500" />
          </div>
          
          <h2 className="text-2xl font-bold mb-3 text-gray-800">File Not Found</h2>
          
          <p className="text-gray-600 mb-8">
            The file you're looking for doesn't exist or has been removed.
          </p>
          
          <Link 
            to="/" 
            className="btn btn-primary inline-flex items-center"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
            </svg>
            Upload New File
          </Link>
        </div>
      </div>
    );
  }

  if (status === 'expired') {
    return (
      <div className="max-w-lg mx-auto text-center py-16 animate-fadeIn">
        <div className="card border-orange-100">
          <div className="rounded-full bg-orange-50 p-4 w-16 h-16 flex items-center justify-center mx-auto mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-orange-500" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
            </svg>
          </div>
          
          <h2 className="text-2xl font-bold mb-3 text-gray-800">File Expired</h2>
          
          <p className="text-gray-600 mb-8">
            This file has expired and is no longer available for download.
            Files on CloudFrog automatically expire after 7 days.
          </p>
          
          <Link 
            to="/" 
            className="btn btn-primary inline-flex items-center"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
            </svg>
            Upload New File
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-lg mx-auto text-center py-16 animate-fadeIn">
      <div className="card border-green-100">
        <div className="rounded-full bg-green-50 p-4 w-16 h-16 flex items-center justify-center mx-auto mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-green-500" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
        </div>
        
        <h2 className="text-2xl font-bold mb-3 text-gray-800">Your download has started</h2>
        
        <p className="text-gray-600 mb-6">
          If your download doesn't begin automatically, click the button below.
        </p>
        
        <a 
          href={`${api.defaults.baseURL}/download/${shortCode}`} 
          className="btn btn-primary inline-flex items-center"
        >
          <ArrowTopRightOnSquareIcon className="h-5 w-5 mr-2" />
          Download File
        </a>
      </div>
    </div>
  );
};

export default DownloadPage;
