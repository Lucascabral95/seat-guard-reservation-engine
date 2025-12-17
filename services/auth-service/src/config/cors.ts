import { envs } from "./envs.schemas";


const MY_URL_FRONTEND = envs.myFrontendUrl;

const corsOptions = {
    origin: [
        'http://localhost:3000',
        'http://localhost:3001',
        'http://localhost:3002',
        'http://localhost:4000',
        'http://localhost:5173',
        'https://doraemon-ecommerce-git-main-lucascabral95s-projects.vercel.app/',
        MY_URL_FRONTEND,
        /^https:\/\/doraemon-ecommerce-.*\.vercel\.app$/,
    ],
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE',],
    allowedHeaders: ['Content-Type', 'Authorization'],
    credentials: true,
};

export default corsOptions;