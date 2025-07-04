import "./App.css";

function App() {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <header className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Web Crawler Dashboard
          </h1>
          <p className="text-gray-600">
            Initial scaffolding - full implementation coming next
          </p>
        </header>

        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-center">
            <div className="animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-3/4 mx-auto mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mx-auto mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-2/3 mx-auto"></div>
            </div>
            <p className="mt-6 text-sm text-gray-500">
              React + TypeScript + Tailwind CSS setup complete
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
