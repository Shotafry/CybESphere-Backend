-- Extensiones necesarias para CybESphere
-- Este script se ejecuta automáticamente al inicializar la base de datos

-- UUID extension para generar IDs únicos
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- PostGIS extension para funcionalidades de geolocalización
-- CREATE EXTENSION IF NOT EXISTS postgis;
-- CREATE EXTENSION IF NOT EXISTS postgis_topology;

-- Full text search en español
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Configuración de texto completo en español
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'spanish_unaccent') THEN
        CREATE TEXT SEARCH CONFIGURATION spanish_unaccent (COPY = spanish);
        ALTER TEXT SEARCH CONFIGURATION spanish_unaccent
            ALTER MAPPING FOR hword, hword_part, word
            WITH unaccent, spanish_stem;
    END IF;
END
$$;

-- Crear función para generar slugs
CREATE OR REPLACE FUNCTION generate_slug(input_text TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN lower(
        regexp_replace(
            unaccent(input_text),
            '[^a-zA-Z0-9]+', '-', 'g'
        )
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Función para calcular distancia usando fórmula de Haversine
-- (Alternativa a PostGIS para cálculos simples de distancia)
CREATE OR REPLACE FUNCTION haversine_distance(
    lat1 DOUBLE PRECISION,
    lon1 DOUBLE PRECISION,
    lat2 DOUBLE PRECISION,
    lon2 DOUBLE PRECISION
) RETURNS DOUBLE PRECISION AS $$
DECLARE
    R CONSTANT DOUBLE PRECISION := 6371; -- Radio de la Tierra en km
    dlat DOUBLE PRECISION;
    dlon DOUBLE PRECISION;
    a DOUBLE PRECISION;
    c DOUBLE PRECISION;
BEGIN
    dlat := radians(lat2 - lat1);
    dlon := radians(lon2 - lon1);
    
    a := sin(dlat/2) * sin(dlat/2) + 
         cos(radians(lat1)) * cos(radians(lat2)) * 
         sin(dlon/2) * sin(dlon/2);
    
    c := 2 * atan2(sqrt(a), sqrt(1-a));
    
    RETURN R * c;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Configuración de timezone por defecto
SET timezone = 'Europe/Madrid';

-- Configuración de búsqueda de texto por defecto
SET default_text_search_config = 'spanish_unaccent';

COMMENT ON DATABASE cybesphere_dev IS 'CybESphere Development Database - Plataforma de Eventos de Ciberseguridad';