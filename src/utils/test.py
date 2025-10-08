import json
import folium

def read_coords_from_json(filename):
    with open(filename, 'r', encoding='utf-8') as file:
        data = json.load(file)

    # Предполагаем, что data - это массив объектов с полем 'geom'
    all_coordinates = []
    for item in data:
        if 'geom' in item and 'coordinates' in item['geom']:
            # Извлекаем координаты и меняем порядок [lat, lon] -> [lon, lat]
            coords = item['geom']['coordinates']
            coords_swapped = [[lon, lat] for lat, lon in coords]
            all_coordinates.append(coords_swapped)
        else:
            print(f"Предупреждение: объект без поля 'geom' или 'coordinates' пропущен: {item}")

    if not all_coordinates:
        raise ValueError("Не найдено ни одного валидного объекта с координатами")
    
    return all_coordinates

# Использование
filename = 'coords.json'
try:
    all_coords_swapped = read_coords_from_json(filename)
    
    # Создаем карту, используя первую точку первого маршрута как центр
    m = folium.Map(location=all_coords_swapped[0][0], zoom_start=13)
    
    # Добавляем каждую линию на карту
    for i, coords in enumerate(all_coords_swapped):
        folium.PolyLine(coords, color="blue", weight=2.5, opacity=1).add_to(m)
        
        # Добавляем маркеры для начала и конца линии (опционально)
        if len(coords) > 0:
            folium.Marker(
                coords[0], 
                popup=f"Начало маршрута {i+1}",
                icon=folium.Icon(color='green', icon='play')
            ).add_to(m)
            
            folium.Marker(
                coords[-1], 
                popup=f"Конец маршрута {i+1}",
                icon=folium.Icon(color='red', icon='stop')
            ).add_to(m)
    
    m.save("route.html")
    print(f"Карта успешно сохранена в route.html. Отображено {len(all_coords_swapped)} маршрутов.")

except FileNotFoundError:
    print(f"Ошибка: Файл {filename} не найден")
except json.JSONDecodeError:
    print(f"Ошибка: Файл {filename} содержит некорректный JSON")
except ValueError as e:
    print(f"Ошибка: {e}")
except Exception as e:
    print(f"Неожиданная ошибка: {e}")