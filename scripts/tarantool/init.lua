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
    --box.space.user_credentials:insert{box.sequence.user_id_sequence:next()}
    --box.space.user_credentials.index.primary:get(1).services['TG']

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

    function add_user_service(user_id, service, login, password)
        local tuple = box.space.user_credentials.index.primary:get{user_id}
        if tuple == nil then
            return box.tuple.new({"User not found", 1})
        end

        local services = tuple[2]
        services[service] = {Login = login, Password = password}

        box.space.user_credentials:update(user_id, {{'=', 2, services}})
        return box.tuple.new({nil, 0})
    end

    box.schema.func.create('add_user_service', {language = 'LUA'})

    function remove_user_service(user_id, service)
        local tuple = box.space.user_credentials.index.primary:get{user_id}
        if tuple == nil then
            return box.tuple.new({"User not found", 2})
        end

        local services = tuple[2]
        if services[service] == nil then
            return box.tuple.new({"Service not found", 1})
        end

        services[service] = nil
        box.space.user_credentials:update(user_id, {{'=', 2, services}})
        return box.tuple.new({nil, 0})
    end

    box.schema.func.create('remove_user_service', {language = 'LUA'})

    function add_user(services)
        local user_id = user_id_sequence:next()
        box.space.user_credentials:insert{user_id, services}
        return user_id
    end

    box.schema.func.create('add_user', {language = 'LUA'})
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
box.schema.user.grant('passman', 'execute', 'function', 'add_user_service')
box.schema.user.grant('passman', 'execute', 'function', 'remove_user_service')
box.schema.user.grant('passman', 'execute', 'function', 'add_user')
box.schema.user.grant('passman', 'read,write,USAGE', 'sequence', 'user_id_sequence')
box.schema.user.grant('passman', 'execute', 'universe')

-- Создаём пользователя для подключения telegram-bot
box.schema.user.create('tgbot', {password='passw0rd', if_not_exists=true})
box.schema.user.grant('tgbot', 'read,write,execute', 'space', 'tg_id')
--require('console').start() os.exit()