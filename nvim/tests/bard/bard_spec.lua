local plugin = require("bard")

describe("setup", function()
  it("works with default", function()
    assert("my first function with param = Hello!", plugin.ask())
  end)

  it("works with custom var", function()
    plugin.setup({ opt = "custom" })
    assert("my first function with param = custom", plugin.ask())
  end)
end)
