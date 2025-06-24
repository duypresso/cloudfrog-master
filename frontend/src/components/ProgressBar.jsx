import React from 'react';
import { motion } from 'framer-motion';

const ProgressBar = ({ progress }) => {
  return (
    <div className="mb-6">
      <div className="flex items-center justify-between mb-1">
        <span className="text-sm font-medium text-blue-600">Uploading...</span>
        <span className="text-sm font-medium text-blue-600">{progress}%</span>
      </div>
      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
        <motion.div
          className="h-full bg-blue-500 rounded-full"
          initial={{ width: 0 }}
          animate={{ width: `${progress}%` }}
          transition={{ duration: 0.3 }}
        />
      </div>
    </div>
  );
};

export default ProgressBar;
