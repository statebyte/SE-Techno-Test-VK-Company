import { useEffect, useState } from 'react';

interface ContainerInfo {
  container_id: string;
  container_name: string;
  ip_address: string;
  status: string;
  updated_at?: string; // может быть undefined, если не приходит
}

function App() {
  const [containers, setContainers] = useState<ContainerInfo[]>([]);
  const [error, setError] = useState<string | null>(null);

  const BASE_API_URL = 'http://localhost:8080';

  useEffect(() => {
    fetchContainers();
  }, []);

  const fetchContainers = async () => {
    setError(null);
    try {
      const response = await fetch(`${BASE_API_URL}/api/containers`);
      if (!response.ok) {
        throw new Error(`Ошибка при запросе: ${response.status}`);
      }
      const data = await response.json();
      setContainers(data);
    } catch (err: any) {
      setError(err.message || 'Неизвестная ошибка');
    }
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: 20 }}>
      <h1>Мониторинг контейнеров</h1>
      {error && <div style={{ color: 'red' }}>Ошибка: {error}</div>}
      <button onClick={fetchContainers}>Обновить</button>

      <table style={{ width: '100%', marginTop: 20, borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ border: '1px solid #ccc', padding: 8 }}>ID</th>
            <th style={{ border: '1px solid #ccc', padding: 8 }}>Имя</th>
            <th style={{ border: '1px solid #ccc', padding: 8 }}>IP</th>
            <th style={{ border: '1px solid #ccc', padding: 8 }}>Статус</th>
            <th style={{ border: '1px solid #ccc', padding: 8 }}>Последнее обновление</th>
          </tr>
        </thead>
        <tbody>
          {containers.map((c) => (
            <tr key={c.container_id}>
              <td style={{ border: '1px solid #ccc', padding: 8 }}>{c.container_id}</td>
              <td style={{ border: '1px solid #ccc', padding: 8 }}>{c.container_name}</td>
              <td style={{ border: '1px solid #ccc', padding: 8 }}>{c.ip_address}</td>
              <td style={{ border: '1px solid #ccc', padding: 8 }}>{c.status}</td>
              <td style={{ border: '1px solid #ccc', padding: 8 }}>
                {c.updated_at ? new Date(c.updated_at).toLocaleString() : '—'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;
