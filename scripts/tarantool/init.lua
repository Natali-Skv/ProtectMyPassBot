box.cfg({listen="127.0.0.1:3301"})

local user_credentials_space = box.space.user_credentials
if not user_credentials_space then
    box.schema.create_space('user_credentials')
    user_credentials_space = box.space.user_credentials
    user_credentials_space:format({
        {name = 'user_id', type = 'unsigned'},
        {name = 'services', type = 'map'},
    })

    user_credentials_space:create_index('primary', {
        type = 'tree',
        parts = {
            {field = 'user_id', type = 'unsigned'},
        },
        unique = true,
    })

    user_id_sequence = box.schema.sequence.create('user_id_sequence')
    box.space.user_credentials:insert{box.sequence.user_id_sequence:next(), {['TG'] = {['Login'] = 'natali_telegram', ['Password'] = 'tg_passw'}, ['VK'] = {['Login'] = 'natali', ['Password'] = 'password'}}}
    --box.space.user_credentials.index.primary:get(1).services['TG']

    --local get_user_credentials = function (user_id, service_key)
    --    local user_credentials_space = box.space.user_credentials
    --    local user_credentials = user_credentials_space.index.primary:get(user_id)
    --    if user_credentials and user_credentials.services[service_key] then
    --        return user_credentials.services[service_key]['Login'], user_credentials.services[service_key]['Password']
    --    end
    --    return nil, nil
    --end

    function get_user_creds (user_id, service_key)
        local user_credentials_space = box.space.user_credentials
        local user_credentials = user_credentials_space.index.primary:get(user_id)
        if user_credentials and user_credentials.services[service_key] then
            local login = user_credentials.services[service_key]['Login']
            local password = user_credentials.services[service_key]['Password']
            t = box.tuple.new({login, password})
            return t
        end
        return box.tuple.new({nil, nil})
    end

    box.schema.func.create('get_user_creds', {language = 'LUA'})


    box.schema.func.create('get_user_credentials', {if_not_exists = true, language = 'LUA', body = get_user_credentials})
end

local tg_id_space = box.space.tg_id
if not tg_id_space then
    tg_id_space = box.schema.create_space('tg_id')
    tg_id_space:format({
        {name = 'tg_id', type = 'unsigned'},
        {name = 'user_id', type = 'unsigned'},
    })
    tg_id_space:create_index('primary', {
        type = 'tree',
        parts = {'tg_id'},
    })
end

-- Создаём пользователя для подключения микросервиса passman
box.schema.user.create('passman', {password='passw0rd', if_not_exists=true})
box.schema.user.grant('passman', 'read,write,execute', 'space', 'user_credentials')
box.schema.user.grant('passman', 'execute', 'function', 'get_user_creds')

-- Создаём пользователя для подключения telegram-bot
box.schema.user.create('tgbot', {password='passw0rd', if_not_exists=true})
box.schema.user.grant('tgbot', 'read,write,execute', 'space', 'tg_id')
--require('console').start() os.exit()